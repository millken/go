package main

///*
import (
	. "./go-gui"
	W "github.com/lxn/go-winapi"
//	"syscall"
)
//*/
func main() {
	rect := RectArea {100,100,500,300}
	hwnd := CreateWindow("这是一个测试","testclass", rect)
	if hwnd == 0 {
		//t.Log("create errot.")
	}

	hwnd1 := CreateWindow("这是另一个窗口","test", rect)
	W.ShowWindow(hwnd, W.SW_SHOW)
	W.ShowWindow(hwnd1, W.SW_SHOW)
	Start(hwnd)
}