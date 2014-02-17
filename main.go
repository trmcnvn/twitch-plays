package main

import (
  "./lib/win32"
  "fmt"
  irc "github.com/fluffle/goirc/client"
  "github.com/laurent22/toml-go"
  "time"
)

type Twitch struct {
  user    string
  token   string
  channel string
}

type Input struct {
  user    string
  message string
  key     uint16
}

var (
  validKeys = map[string]uint16{
    "a":     win32.VK_X,
    "b":     win32.VK_Z,
    "x":     win32.VK_S,
    "y":     win32.VK_A,
    "up":    win32.VK_UP,
    "left":  win32.VK_LEFT,
    "down":  win32.VK_DOWN,
    "right": win32.VK_RIGHT,
  }
)

func inputHandler(in <-chan Input) {
  for {
    input := <-in
    fmt.Printf("Handling: %v - %v\n", input.user, input.message)

    win32.SendInput(win32.INPUT{
      Type: win32.INPUT_KEYBOARD,
      Ki: win32.KEYBDINPUT{
        WVk:     input.key,
        DwFlags: win32.KEYEVENTF_KEYDOWN,
      },
    })

    time.Sleep(200 * time.Millisecond)

    win32.SendInput(win32.INPUT{
      Type: win32.INPUT_KEYBOARD,
      Ki: win32.KEYBDINPUT{
        WVk:     input.key,
        DwFlags: win32.KEYEVENTF_KEYUP,
      },
    })
  }
}

func main() {
  var p toml.Parser
  d := p.ParseFile("config/app.conf")
  twitch := Twitch{
    user:    d.GetString("twitch.user"),
    token:   d.GetString("twitch.token"),
    channel: "#" + d.GetString("twitch.channel"),
  }

  in := make(chan Input)
  go inputHandler(in)

  // make sure emulator window is open
  h := win32.FindWindow("DeSmuME", "DeSmuME 0.9.10 x64")
  if h == nil {
    panic("Couldn't find emulator window.")
  }
  fmt.Printf("Emulator Window: %v\n", h)

  // connect to TwitchTV chat
  c := irc.SimpleClient(twitch.user)

  c.AddHandler("connected", func(conn *irc.Conn, line *irc.Line) {
    c.Join(twitch.channel)
  })

  quit := make(chan bool)
  c.AddHandler("disconnected", func(conn *irc.Conn, line *irc.Line) {
    quit <- true
  })

  c.AddHandler("privmsg", func(conn *irc.Conn, line *irc.Line) {
    fmt.Println(line.Raw)

    if line.Args[0] != twitch.channel {
      return
    }

    user, message := line.Nick, line.Args[1]
    key, ok := validKeys[message]
    if !ok {
      return
    }

    in <- Input{
      user:    user,
      message: message,
      key:     key,
    }
  })

  if err := c.Connect("irc.twitch.tv", twitch.token); err != nil {
    panic("Couldn't connect to TwitchTV.")
  }

  // keep the connection alive
  <-quit
}
