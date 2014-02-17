package win32

type INPUT struct {
  Type uint32
  Ki   KEYBDINPUT
}

type KEYBDINPUT struct {
  WVk         uint16
  WScan       uint16
  DwFlags     uint32
  Time        uint32
  DwExtraInfo uintptr
}