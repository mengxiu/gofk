package utils

import (
	"fmt"
	"github.com/lxn/win"
	"io/ioutil"
	"log"
	"path/filepath"
	"syscall"
)

func SwitchWindow(handle syscall.Handle) {
	var a, b win.RECT
	hWnd := win.HWND(handle)
	nCmdShow := int32(9)
	win.ShowWindow(hWnd, nCmdShow)
	win.SetForegroundWindow(hWnd)
	win.GetClientRect(hWnd, &a)
	win.GetWindowRect(hWnd, &b)

}

type folder map[string]interface{}

// ListDir 传入绝对路径,将文件夹映射成map，但暂时没想到好的解析方式
func ListDir(rootPath string) folder {
	//type folder map[string]interface{}

	var Fold folder
	Fold = make(folder)
	files, errDir := ioutil.ReadDir(rootPath)
	if errDir != nil {
		log.Fatal(errDir)
	}

	for _, file := range files {
		if file.IsDir() {
			Fold[file.Name()] = ListDir(rootPath + "/" + file.Name())

		} else {

			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(rootPath + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}
			Fold[file.Name()] = strAbsPath

		}
	}
	return Fold
}

// FindFileBySuffix 传入根目录，根据后缀名返回文件的绝对路径列表,
//传入的后缀名要带”.“
func FindFileBySuffix(dir, suffix string) (fileList []string) {
	path, err := filepath.Abs(dir)
	fmt.Println(path)
	if err != nil {
		log.Fatal(err)
	}
	files, errDir := ioutil.ReadDir(path)
	if errDir != nil {
		log.Fatal(errDir)
	}
	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
			if (filepath.Ext(file.Name())) == suffix {

				fileList = append(fileList, path+"\\"+file.Name())
			}
		}
	}
	return
}
