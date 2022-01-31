package Account

import (
	"GoFK/LockAndHook"
	"GoFK/utils"
	"fmt"
	"github.com/lxn/win"
	"gopkg.in/ini.v1"
	"time"
)

var defaultConfig Config

func init() {
	Config, err := ini.Load("config.ini")
	if err != nil {
		println("读取配置文件失败，请检查同级目录下的config.ini文件")
		panic(err)
	}
	err = Config.MapTo(&defaultConfig)
	if err != nil {
		fmt.Println("映射配置文件到结构体出错:" + err.Error())
	}
}

type ACCOUNT struct {
	hwnd    win.HWND
	cRect   win.RECT
	wRect   win.RECT
	Gui     utils.GUI
	srcPath string
}
type Config struct {
	FilePath   string `ini:"filepath" json:"filepath,omitempty" `
	Method     string `ini:"method" json:"method,omitempty" `
	IsTeamMode string `ini:"is_team_mode" json:"is_team_mode,omitempty" `
	WinTitle   string `ini:"win_title"`
}

type Wins struct {
	accountList []ACCOUNT
	config      Config
}

//置顶窗口，并检测位置变动
func (self *ACCOUNT) AtivateWin() {
	<-LockAndHook.GuiLock
	//ShowWindows很耗费时间，50-80ns，但是配合参数win.SW_RESTORE=9能把最小化的窗口还原
	//此外win.ShowWindow还需要管理员权限运行
	win.ShowWindow(self.hwnd, win.SW_RESTORE)
	//win.SetForegroundWindow耗费6-8ns
	win.SetForegroundWindow(self.hwnd)
	wRect := win.RECT{}
	win.GetWindowRect(self.hwnd, &(wRect))
	println("偏移", wRect.Left-self.wRect.Left, wRect.Top-self.wRect.Top)
	self.Gui.Translate(wRect.Left-self.wRect.Left, wRect.Top-self.wRect.Top)
	self.wRect = wRect
	self.Gui.Init(&self.wRect)
	LockAndHook.GuiLock<-1
}
func (self *ACCOUNT) Init(handle win.HWND) {
	self.hwnd = handle
	fmt.Println("handle为", handle)
	win.ShowWindow(self.hwnd, win.SW_RESTORE)
	//win.SetForegroundWindow耗费6-8ns
	win.SetForegroundWindow(self.hwnd)
	time.Sleep(time.Second * 1)
	//t1 := time.Now()
	win.GetWindowRect(self.hwnd, &(self.wRect))
	win.GetClientRect(self.hwnd, &(self.cRect))
	self.Gui.Init(&self.wRect)
	//fmt.Println(self.wRect, self.cRect, w, c)
	//robotgo.MoveSmooth(int(self.wRect.Right-self.cRect.Right),int(self.wRect.Bottom-self.cRect.Bottom))
	//fmt.Println(robotgo.GetMousePos())
	//robotgo.MoveSmooth(int(self.wRect.Left),int(self.wRect.Top))
	//fmt.Println(time.Since(t1) / time.Millisecond)
}
func (self *ACCOUNT) SetSrcPath(path string) {
	self.srcPath = path
}

func (self *Wins) Start() {
	self.config = defaultConfig
	handleList, err := utils.FindWindows(self.config.WinTitle)
	if err != nil {
		fmt.Println(err)
		return
	}
	//使用for range一定要再三小心谨慎
	for _, handle := range handleList {
		acc := ACCOUNT{}
		acc.Init(win.HWND(handle))
		self.accountList = append(self.accountList, acc)
	}

	if self.config.IsTeamMode == "true" {
		self.teamMode()
	} else {
		self.soloMode()
	}
}
