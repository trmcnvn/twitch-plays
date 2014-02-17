package win32

import (
  "syscall"
  "unsafe"
)

var (
  moduser32 = syscall.NewLazyDLL("user32.dll")
  procFindWindow = moduser32.NewProc("FindWindowW")
  procSendInput = moduser32.NewProc("SendInput")
)

func SendInput(input INPUT) uint {
  var inputArr []INPUT
  inputArr = append(inputArr, input)

  // 0x28 is the correct size of C.INPUT
  ret, _, _ := procSendInput.Call(uintptr(len(inputArr)), uintptr(unsafe.Pointer(&inputArr[0])), uintptr(0x28))
  return uint(ret)
}

func FindWindow(className string, windowName string) unsafe.Pointer {
  lpClassName := syscall.StringToUTF16Ptr(className)
  lpWindowName := syscall.StringToUTF16Ptr(windowName)

  ret, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(lpClassName)), uintptr(unsafe.Pointer(lpWindowName)))
  return unsafe.Pointer(ret)
}