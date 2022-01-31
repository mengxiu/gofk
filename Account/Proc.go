package Account

import (
	"GoFK/LockAndHook"
	"GoFK/utils"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	guiLock chan int
)

type fold map[string]interface{}

func doProcWithGuiLock(proc func()) {
	guiLock <- 1
	proc()
	<-guiLock
}

type MethodMapsType map[string]reflect.Value

var MethodMaps MethodMapsType

func init() {
	MethodMaps = make(MethodMapsType, 0)
	account := ACCOUNT{}
	value := reflect.ValueOf(&account)
	vft := value.Type()
	mNum := value.NumMethod()

	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		MethodMaps[mName] = value.Method(i)

	}

}
func (self *Wins) soloMode() {
	wg := sync.WaitGroup{}
	for account := range self.accountList {
		wg.Add(1)
		println(account)
	}

}
func (self *Wins) teamMode() {

	vtf := reflect.TypeOf(self)
	proc, _ := vtf.MethodByName(self.config.Method)
	proc.Func.Call([]reflect.Value{reflect.ValueOf(self)})
}

func (self *Wins) NormalProc() error{
	rand.Seed(time.Now().UnixNano())
	exConfig, err := ini.Load(self.fWFP("config.ini"))
	fmt.Println(exConfig.SectionStrings())
	if err != nil {
		fmt.Println("图片文件夹的配置文件读取出错", err.Error())
		return err
	}
	//对应文件夹中的图片的绝对路径的列表
	var start, beforeEnd, over, confirm,xuanshang,yuhunMan,jiacheng []string
	//配置参数
	var (
		mainRest          = exConfig.Section("main").Key("rest").MustBool()
		mainTime          = exConfig.Section("main").Key("time").MustInt()
		jiachengLong      = exConfig.Section("jiacheng").Key("long").MustInt()
		jiachengShort     = exConfig.Section("jiacheng").Key("short").MustInt()
		jiachengClose     = exConfig.Section("jiacheng").Key("close").MustBool()
		otherShutdownGame = exConfig.Section("other").Key("shutdown_game").MustBool()
		otherShutdown     = exConfig.Section("other").Key("shutdown").MustBool()
	)
	
	start = utils.FindFileBySuffix(self.config.FilePath+"/start", ".png")
	beforeEnd = utils.FindFileBySuffix(self.config.FilePath+"/before_end", ".png")
	over = utils.FindFileBySuffix(self.config.FilePath+"/over", ".png")
	confirm = utils.FindFileBySuffix(self.config.FilePath+"/confirm", ".png")
	xuanshang=utils.FindFileBySuffix("src/exception/xuanshang",".png")
	yuhunMan=utils.FindFileBySuffix("src/exception/yuhunMan",".png")
	jiacheng=utils.FindFileBySuffix("src/exception/jiacheng",".png")
//fmt.Println("查看",start, beforeEnd, over, confirm,mainTime,mainRest,jiachengLong,jiachengShort,jiachengClose,otherShutdownGame,otherShutdown)

	flag := false
	count := 0
	for {
		//start
		for j, acc := range self.accountList {
			acc.AtivateWin()
			startFlag:=0
			for  {
				println("start",j)
				_,_,ok:=acc.Gui.FindImageAndRecord(start[rand.Intn(len(start))])
				if ok{
					LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,start[rand.Intn(len(start))],15,2)
					robotgo.Click()
					startFlag++
				}else {
					break
				}
				time.Sleep(time.Millisecond*200)
				if startFlag>15{
					return errors.New("无法开始")
				}

			}

		}

		time.Sleep(time.Second*time.Duration(mainTime))

		//beforeend
		if len(beforeEnd)>0{
			for _, acc := range self.accountList {
				acc.AtivateWin()
				startFlag:=0
				for  {
					if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,beforeEnd[rand.Intn(len(beforeEnd))],15,2){
						robotgo.Click()
					}
					_,_,ok:=acc.Gui.FindImageAndRecord(over[rand.Intn(len(over))])
					if ok{
						break
					}
					startFlag++
					time.Sleep(time.Millisecond*200)
					if startFlag>60{
						break
					}

				}

			}
		}
		//over，confirm，exception
		for j, acc := range self.accountList {
			acc.AtivateWin()
			startFlag:=0
			for  {
				println("over",j)
				if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,over[rand.Intn(len(over))],15,2){
					robotgo.Click()
				}
				if _,_,ok:=acc.Gui.FindImageAndRecord(confirm[rand.Intn(len(confirm))]);ok{
					break
				}

				startFlag++
				time.Sleep(time.Millisecond*200)
				if startFlag>30{
					if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,xuanshang[rand.Intn(len(xuanshang))],15,2){
						robotgo.Click()
					}
					if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,yuhunMan[rand.Intn(len(yuhunMan))],15,2){
						robotgo.Click()
						flag=true
						time.Sleep(time.Millisecond*500)
						if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,over[rand.Intn(len(over))],15,2){
							robotgo.Click()
						}
						break
					}

				}

			}

		}

		//for i:=0;i<len(self.accountList);i++{
		//	self.accountList[i].AtivateWin()
		//	LockAndHook.FuncWithLock(self.accountList[i].Gui.FindAndMoveTo,self.fWFP("start\\test3.png"),1,2)
		//	//self.accountList[i].Gui.FindAndMoveTo(self.fWFP("start\\test3.png"),1,2)
		//}


		if flag {
			if jiachengClose{
				for _, acc := range self.accountList{
					if LockAndHook.FuncWithLock(acc.Gui.FindAndMoveTo,jiacheng[rand.Intn(len(jiacheng))],15,2) {
						robotgo.MoveSmoothRelative(0,jiachengShort+rand.Intn(jiachengLong-jiachengShort))
						robotgo.Click()
					}
				}
			}


			break
		}
		if mainRest && count-30 > rand.Intn(10) {
			count -= 30

			time.Sleep(time.Second * time.Duration(60+rand.Intn(30)))
		}
	}
	if otherShutdownGame{
		cmd:="taskkill /f /t /im onmyoji.exe"
		_, err := exec.Command("cmd", "/c", cmd).Output()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if otherShutdown{
		cmd:="shutdown -s -t 1"
		_, err := exec.Command("cmd", "/c", cmd).Output()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	println("over")
	return nil
}

//返回self.config.FilePath加传入字符串的拼接
func (self *Wins) fWFP(str string) string {
	buffer := strings.Builder{}
	defer buffer.Reset()
	abspath, _ := filepath.Abs(self.config.FilePath)
	buffer.WriteString(abspath)
	buffer.WriteString("\\")
	buffer.WriteString(str)
	return buffer.String()
}

//传入图片的绝对路径
func (self *ACCOUNT) clickAndDetectlAppear(imgToClick, imgToAppear string) bool {

	if self.Gui.FindAndMoveTo(imgToClick) {
		robotgo.Click()
	}

	time.Sleep(time.Millisecond * 300)
	_, _, flag2 := self.Gui.FindImageAndRecord(imgToAppear)

	return flag2
}
func (self *ACCOUNT) clickAndDetectDisappear(imgToClick string) bool {
	flag := self.Gui.FindAndMoveTo(imgToClick)
	if flag {
		robotgo.Click()
	}
	return flag
}

func ListDir(folder string) {
	files, errDir := ioutil.ReadDir(folder)
	if errDir != nil {
		log.Fatal(errDir)
	}

	for _, file := range files {
		if file.IsDir() {
			ListDir(folder + "/" + file.Name())
		} else {
			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(folder + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}

			fmt.Println(strAbsPath)
		}
	}

}
