//https://gist.github.com/thesubtlety/be6e7ec9c19083473bed4cae11c8160d
//https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32Dll      = windows.NewLazyDLL("user32.dll")
	procCursorInfo = user32Dll.NewProc("GetCursorPos")
)

func main() {

	type Point struct {
		X, Y int32
	}

	var point Point
	for 1 > 0 {
		procCursorInfo.Call(uintptr(unsafe.Pointer(&point)))
		fmt.Println(point)
	}
}
