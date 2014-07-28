package main

import (
	"syscall"
	"unsafe"
	"strconv"
	"fmt"
)

type OSVERSIONINFO struct {
	dwOSVersionInfoSize int32
	dwMajorVersion int32
	dwMinorVersion int32
	dwBuildNumber int32
	dwPlatformId int32
	szCSDVersion [128]byte
}

func main(){
	fmt.Println("操作系统版本" + GetOsVersion());

}

func IntPtr(n int) uintptr {
	return uintptr(n);
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)));
}

func ShowMessage(title string,text string) {
	user32 := syscall.NewLazyDLL("user32.dll");
	MessageBoxW :=  user32.NewProc("MessageBoxW");
	MessageBoxW.Call(IntPtr(0),StrPtr(text),StrPtr(title),IntPtr(0));
}

func GetOsVersion() string {
	var version string = "Unknown Version";
	kernel32 := syscall.NewLazyDLL("kernel32.dll");
	var os OSVERSIONINFO;
	os.dwOSVersionInfoSize = int32(unsafe.Sizeof(os));
	GetVersionExA := kernel32.NewProc("GetVersionExA");
	rt,_,_ := GetVersionExA.Call(uintptr(unsafe.Pointer(&os)));
	if int(rt)==1 {
		switch {
			case os.dwMajorVersion==4 && os.dwMinorVersion==0 && os.dwPlatformId==1:
				version = "Windows 95";
				break;
			case os.dwMajorVersion==4 && os.dwMinorVersion==10:
			    version = "Windows 98";
				break;
			case os.dwMajorVersion==4 && os.dwMinorVersion==90:
				version = "Windows Me";
				break;
			case os.dwMajorVersion==4 && os.dwMinorVersion==0 && os.dwPlatformId==2:
				version = "Windows NT4";
				break;
			case os.dwMajorVersion==5 && os.dwMinorVersion==0:
				version = "Windows 2000";
				break;
			case os.dwMajorVersion==5 && os.dwMinorVersion==1:
				version = "Windows XP";
				break;
			case os.dwMajorVersion==5 && os.dwMinorVersion==2:
				version = "Windows 2003";
				break;
			case os.dwMajorVersion==6 && os.dwMinorVersion==0:
				version = "Windows Vista";
				break;
			case os.dwMajorVersion==7 && os.dwMinorVersion==0:
				version = "Windows 7";
				break;
			default:
				version = "Windows 8";
				break;
		}
		version = version + " Build("+strconv.FormatInt(int64(os.dwBuildNumber),10)+") " + string(os.szCSDVersion[0:]);
	}
	return version;
}
