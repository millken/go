//http://www.du52.com/text.php?id=54
package main

import "syscall"
import "unsafe"

func main() {
    var hand uintptr = uintptr(0);
    var operator uintptr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("open")));
    var fpath uintptr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("D:\Program Files\TTPlayer\TTPlayer.exe")));
    var param uintptr = uintptr(0);
    var dirpath uintptr = uintptr(0);
    var ncmd uintptr = uintptr(1);
    shell32 := syscall.NewLazyDLL("shell32.dll");
    ShellExecuteW := shell32.NewProc("ShellExecuteW");
    _,_,_ = ShellExecuteW.Call(hand,operator,fpath,param,dirpath,ncmd);
}
