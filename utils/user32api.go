package utils

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	user32api          = windows.NewLazySystemDLL("User32.dll")
	procEnumWindows    = user32api.NewProc("EnumWindows")
	procGetWindowTextW = user32api.NewProc("GetWindowTextW")
)

func StringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0)
	return &chars[0]
}

// 回调函数，用于EnumWindows中的回调函数，第一个参数是hWnd，第二个是自定义穿的参数
func AddElementFunc(hWnd syscall.Handle, hWndList *[]syscall.Handle) uintptr {
	*hWndList = append(*hWndList, hWnd)
	return 1
}

// 获取桌面下的所有窗口句柄，包括没有Windows标题的或者是窗口的。
func GetDesktopWindowHWND() []syscall.Handle {
	var hWndList []syscall.Handle
	hL := &hWndList
	r1, _, err := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(syscall.NewCallback(AddElementFunc)), uintptr(unsafe.Pointer(hL)), 0)
	if err != 0 {
		fmt.Println(err)
	}
	fmt.Println(r1)
	//fmt.Println(hWndList)
	return hWndList
}

//枚举所有窗口的句柄
func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//根据句柄获取窗口标题
func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//搜索特定标题的窗口，返回符合条件的窗口句柄的列表
func FindWindows(title string) ([]syscall.Handle, error) {
	var hWndList []syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hWndList = append(hWndList, h)
			return 1 // stop enumeration 已修改此处以实现获取多个阴阳师窗口，原本返回0
		}
		return 1 // continue enumeration
	})
	err := EnumWindows(cb, 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	if hWndList == nil {
		return hWndList, fmt.Errorf("No window with title '%s' found", title)
	}
	fmt.Println(hWndList)
	return hWndList, nil
}
