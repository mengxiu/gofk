package LockAndHook

import (
	"github.com/go-vgo/robotgo"
)

//鼠标，窗口的控制权锁，防止线程切换时窗口和鼠标的冲突
var (
	GuiLock      chan int
	exChangeFlag bool = true
)

func init() {
	GuiLock = make(chan int, 1)
	GuiLock <- 1
}

//全局锁
func GlobalLock() {
	KeyboardHook(exChange, rune(123))
	go MouseSecure()
}
func MouseSecure() {
	for {
		x, y := robotgo.GetMousePos()
		if x*x+y*y < 32 {
			panic("mouseSecure")
		}
	}
}
//暂时没有想到除了反射之外的传入任意函数的修饰方法，也不想用封装成interface,
//实在是太麻烦了，此处仅实现了特定函数类型的装饰器
func FuncWithLock(work func(imagePath string, mouseSpeed ...int) bool,imagePath string, mouseSpeed ...int) bool{
	<-GuiLock

	flag:=work(imagePath,mouseSpeed...)
	GuiLock <- 1
	return flag
}
func exChange() {
	if exChangeFlag {
		println("f12被按下，正在暂停")
		exChangeFlag=!exChangeFlag
		<-GuiLock
	} else {
		println("f12被按下，正在解除暂停")
		exChangeFlag=!exChangeFlag
		GuiLock <- 1
	}
}
