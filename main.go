package main

import (
  "./lib/win32"
  "fmt"
  "github.com/laurent22/toml-go"
  "github.com/thoj/go-ircevent"
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
  // emulator key: windows virtual key
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

func inputHandler(inputs <-chan Input) {
  for {
    input := <-inputs
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
  var parser toml.Parser
  doc := parser.ParseFile("config/app.conf")
  twitch := Twitch{
    user:    doc.GetString("twitch.user"),
    token:   doc.GetString("twitch.token"),
    channel: "#" + doc.GetString("twitch.channel"),
  }

  inputs := make(chan Input)
  go inputHandler(inputs)

  // make sure emulator window is open
  handle := win32.FindWindow("DeSmuME", "DeSmuME 0.9.10 x64")
  if handle == nil {
    panic("Couldn't find emulator window.")
  }
  fmt.Printf("Emulator Window: %v\n", handle)

  // connect to TwitchTV chat
  conn := irc.IRC(twitch.user, twitch.user)
  conn.Password = twitch.token
  err := conn.Connect("irc.twitch.tv:6667")
  if err != nil {
    panic("Couldn't connect to twitch.")
  }

  conn.AddCallback("001", func(e *irc.Event) {
    conn.Join(twitch.channel)
    fmt.Printf("Joined the channel: %v\n", twitch.channel)
  })

  conn.AddCallback("PRIVMSG", func(e *irc.Event) {
    if e.Arguments[0] != twitch.channel {
      return
    }

    user, message := e.Nick, e.Message()
    key, ok := validKeys[message]
    if !ok {
      return
    }

    inputs <- Input{
      user:    user,
      message: message,
      key:     key,
    }
  })

  // keep the connection alive
  conn.Loop()
}
