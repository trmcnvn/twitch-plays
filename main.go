package main

import (
  "fmt"
  "unsafe"
  "time"
  "github.com/laurent22/toml-go"
  "github.com/thoj/go-ircevent"
  "./lib"
)

var (
  // emulator key: windows virtual key
  validKeys = map[string]uint16{
    "a": win32.VK_X,
    "b": win32.VK_Z,
    "x": win32.VK_S,
    "y": win32.VK_A,
    "up": win32.VK_UP,
    "left": win32.VK_LEFT,
    "down": win32.VK_DOWN,
    "right": win32.VK_RIGHT,
  }
)

// makes sure our config file includes the values we need
func validateConf(doc toml.Document) {
  _, user := doc.GetValue("twitch.user")
  _, token := doc.GetValue("twitch.token")
  _, channel := doc.GetValue("twitch.channel")

  if !user || !token || !channel {
    panic("Couldn't validate config file.")
  }
}

// handles messages sent to the twitch channel
func handleMessage(handle unsafe.Pointer, e *irc.Event) {
  user, message := e.Nick, e.Message()
  key, ok := validKeys[message]
  if ok {
    // send the input as a keydown event
    win32.SendInput(win32.INPUT{
      Type: win32.INPUT_KEYBOARD,
      Ki: win32.KEYBDINPUT{
        WVk: key,
        WScan: 0,
        DwFlags: win32.KEYEVENTF_KEYDOWN,
        Time: 0,
        DwExtraInfo: 0,
      },
    })

    // give it some time
    time.Sleep(200 * time.Millisecond)

    // resend the key as a keyup event
    win32.SendInput(win32.INPUT{
      Type: win32.INPUT_KEYBOARD, 
      Ki: win32.KEYBDINPUT{
        WVk: key,
        WScan: 0,
        DwFlags: win32.KEYEVENTF_KEYUP,
        Time: 0,
        DwExtraInfo: 0,
      },
    })

    // output on our console
    fmt.Printf("%v: %v\n", user, message)
  }
}
func
 main() {
  // parse and validate config
  var parser toml.Parser
  doc := parser.ParseFile("config/app.conf")
  validateConf(doc)

  // find emulator window
  handle := win32.FindWindow("DeSmuME", "DeSmuME 0.9.10 x64")
  //handle := win32.FindWindow("Notepad", "Untitled - Notepad")
  if handle == nil {
    panic("Couldn't find emulator window.")
  }
  fmt.Printf("Emulator Window: %v\n", handle)

  // connect to twitch IRC
  user := doc.GetString("twitch.user")
  conn := irc.IRC(user, user)
  conn.Password = doc.GetString("twitch.token")
  err := conn.Connect("irc.twitch.tv:6667")
  if err != nil {
    panic("Couldn't connect to twitch.")
  }

  // join channel on connectionc
  channel := "#" + doc.GetString("twitch.channel")
  conn.AddCallback("001", func (e *irc.Event) {
    conn.Join(channel)
    fmt.Printf("Joined the channel: %v\n", channel)
  })

  // handle messages to the channel
  conn.AddCallback("PRIVMSG", func(e *irc.Event) {
    if e.Arguments[0] == channel {
      handleMessage(handle, e)
    }
  })

  // keep the connection alive
  conn.Loop()
}