package main

import (
	"GoFK/Account"
	"GoFK/LockAndHook"
	"GoFK/utils"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"time"
)

func main() {
	go LockAndHook.GlobalLock()
	//go mouseSecure()
	t1 := time.Now()
	//fmt.Println("a")
	test2()
	//a := utils.GetDesktopWindowHWND()

	//fmt.Println(b)
	//for _, v := range a {
	//	//fmt.Println(i)
	//	if len(b)>0&&int(v) == int(b[0]) {
	//		fmt.Println("True")
	//	}
	//}

	//b, _ := utils.FindWindows("阴阳师-网易游戏")
	//var rect,rect2 win.RECT
	//fmt.Println(robotgo.GetScreenSize())
	//win.GetClientRect(win.HWND(b[0]),&rect2)
	//win.GetWindowRect(win.HWND(b[0]),&rect)
	//fmt.Println(rect,rect2)
	//win.SetForegroundWindow(win.HWND(b[0]))
	//
	//time.Sleep(time.Second*1)

	//bitmap:=robotgo.CaptureScreen(int(float64(rect.Left)*1.5),int(float64(rect.Top)*1.5),int(float64(rect.Right-rect.Left)*1.5),int(float64(rect.Bottom-rect.Top)*1.5))
	//	//bitmap:=robotgo.CaptureScreen(int(float64(rect.Left)),int(float64(rect.Top)),int(float64(rect.Right-rect.Left)),int(float64(rect.Bottom-rect.Top)))
	//fmt.Println(int(rect.Left),int(rect.Top),int(rect.Right-rect.Left),int(rect.Bottom-rect.Top))
	//robotgo.SaveBitmap(bitmap,"test.tif")
	//robotgo.FreeBitmap(bitmap)

	//robotgo.MoveSmooth(277,129)
	//tes:=Account.Account{hwnd: win.HWND(b[0])}
	//tes.Init()

	//robotgo.Move(0,0)
	//robotgo.MoveSmooth(640,640,0.5,1.5)

	//	abfg:=utils.GUI{}
	//	abfg.MoveTo(0,0)

	//fmt.Println(utils.Tweening(0.1,"EaseInOutBounce"))
	t2 := time.Since(t1) / time.Millisecond
	fmt.Println(t2)
	time.Sleep(time.Second * 15)
}
func MouseSecure() {
	for {
		x, y := robotgo.GetMousePos()
		if x*x+y*y < 32 {
			panic("mouseSecure")
		}
	}
}
func test2() {
	wins := Account.Wins{}
	wins.Start()
}
func test() {
	handleList, _ := utils.FindWindows("阴阳师-网易游戏")
	if len(handleList) == 0 {
		return
	}
	accs := make([]Account.ACCOUNT, 0)
	for _, handle := range handleList {
		acc := Account.ACCOUNT{}
		acc.Init(win.HWND(handle))
		accs = append(accs, acc)
	}
	println(len(accs))
	for I := 0; I < 60; I++ {
		for _, acc := range accs {
			println(I)
			acc.AtivateWin()
			acc.Gui.FindAndMoveTo("test3.png", 5)
			acc.Gui.FindAndMoveTo("test4.png", 5)
			acc.Gui.FindAndMoveTo("test3.png", 5)
			acc.Gui.FindAndMoveTo("test2.png", 5)
			acc.Gui.FindAndMoveTo("test4.png", 5)
			acc.Gui.FindAndMoveTo("test2.png", 5)

			time.Sleep(time.Second * 5)
		}

	}

	var rect, rect2 win.RECT
	fmt.Println(robotgo.GetScreenSize())
	win.GetClientRect(win.HWND(handleList[0]), &rect2)
	win.GetWindowRect(win.HWND(handleList[0]), &rect)
	fmt.Println(rect, rect2)
	win.SetForegroundWindow(win.HWND(handleList[0]))
	time.Sleep(time.Second * 1)
}
