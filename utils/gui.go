package utils

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/pkg/errors"
	"gocv.io/x/gocv"
	"image"
	"math"
	"math/rand"
	"path/filepath"
	"time"
)

var (
	dpi float64
)

func init() {
	rand.Seed(time.Now().UnixNano())
	_hwnd := win.GetDesktopWindow()
	hdc := win.GetDC(_hwnd)
	dpi = float64(win.GetDeviceCaps(hdc, win.DESKTOPHORZRES)) / float64(win.GetSystemMetrics(win.SM_CXSCREEN))
	println("系统屏幕缩放比例(dpi):", dpi)
}

type region struct {
	Left   int
	Top    int
	Width  int
	Height int
}

//返回一个向外扩充n个像素的region,n可以是负数
func (self *region) expand(n int) region {
	return region{
		Left:   self.Left - n,
		Top:    self.Top - n,
		Width:  self.Width + n*2,
		Height: self.Height + n*2,
	}
}
func (self *region) multiply(dpi float64) region {
	return region{
		Left:   int(float64(self.Left) * dpi),
		Top:    int(float64(self.Top) * dpi),
		Width:  int(float64(self.Width) * dpi),
		Height: int(float64(self.Height) * dpi),
	}
}

func (self *region) translate(x, y int) region {
	return region{
		Left:   self.Left + x,
		Top:    self.Top + y,
		Width:  self.Width + x,
		Height: self.Height + y,
	}
}

type picPosition map[string]region

type GUI struct {
	//由系统给出的在屏幕上的位置
	wRect *win.RECT
	//经过缩放之后在画面像素中的位置
	picPositionsMap picPosition
	winRegion       region
}

func (self *GUI) Translate(x, y int32) {
	//此处注意range对map操作的经典陷阱,注释中为错误处理
	//for _,v:=range self.picPositionsMap{
	//	v=v.translate(int(float64(x)*dpi),int(float64(y)*dpi))
	//}
	for i, v := range self.picPositionsMap {
		self.picPositionsMap[i] = v.translate(int(float64(x)*dpi), int(float64(y)*dpi))
	}
}

// FindImageAndRecord 定位图片在窗口中的位置，并记录下检测过的图片的具体位置,返回模板在屏幕像素中的具体位置
//大图中没有找到要匹配的图片时返回空区域和False
func (self *GUI) FindImageAndRecord(imagePath string) (region, error, bool) {
	absPath, _ := filepath.Abs(imagePath)

	image1 := gocv.IMRead(absPath, gocv.IMReadColor)
	if image1.Rows()==0{
		return region{}, errors.New("img not found"), false
	}
	defer func(image1 *gocv.Mat) {
		err := image1.Close()
		if err != nil {
			println(err)
		}
	}(&image1)
	if regionTmp, ok := self.picPositionsMap[absPath]; ok != true {
		fmt.Println("not found in map", imagePath,self.picPositionsMap)
		time.Sleep(time.Millisecond * 100)
		//self.MoveTo(int(math.Max(float64(self.winRegion.Left)-100, 15)), int(math.Max(float64(self.winRegion.Top)-100, 15)), 5,1)
		self.MoveTo(int(self.wRect.Left+100), int(self.wRect.Bottom-100), 10, 1)
		image2, err := self.winScreenShot(self.winRegion)

		if err != nil {
			panic(err)
		}
		startPoint, ok := self.matchLocate(image2, image1)
		if ok != true {
			return region{}, errors.New("img not found"), false
		}
		//endPoint:=image.Point{startPoint.X+image1.Rows(),startPoint.Y+image1.Cols()}
		//if self.picPositionsMap == nil {
		//	fmt.Println("new")
		//	self.picPositionsMap = make(picPosition, 0)
		//}
		self.picPositionsMap[absPath] = region{
			Left:   self.winRegion.Left + startPoint.X,
			Top:    self.winRegion.Top + startPoint.Y,
			Width:  image1.Cols(),
			Height: image1.Rows(),
		}
		_,ok2:=self.picPositionsMap[absPath]
		fmt.Println(ok2,"存在测试",self.picPositionsMap)
		return self.picPositionsMap[absPath], nil, true
	} else {
		regionTmp2 := regionTmp.expand(10)
		//fmt.Println(regionTmp.Left, regionTmp.Top)

		time.Sleep(time.Millisecond * 100)
		fmt.Println("ok")
		self.MoveTo(int(math.Max(float64(regionTmp.Left)-150, 15)), int(math.Max(float64(regionTmp.Left)-150, 15)), 5, 1)
		//self.MoveTo(int(self.wRect.Left),int(self.wRect.Bottom),10,1)
		image2, _ := self.winScreenShot(regionTmp2)
		_, ok := self.matchLocate(image2, image1)
		if ok {
			return regionTmp, nil, true
		}
	}

	return region{}, errors.New("img not found"), false
}

// FindAndMoveTo 寻找图片位置并移动到该图片区域的随机一点
//可选参数数量为一个，表示鼠标移动速度，范围1-15,默认值为3
func (self *GUI) FindAndMoveTo(imagePath string, mouseSpeed ...int) bool {
	region0, _, ok := self.FindImageAndRecord(imagePath)
	if ok != true {
		fmt.Println("error")

		return false
	}
	region1 := region0.multiply(1 / dpi)
	randomPoint := image.Point{
		X: rand.Intn(region1.Width) + region1.Left,
		Y: rand.Intn(region1.Height) + region1.Top,
	}
	if len(mouseSpeed) > 0 {
		self.MoveTo(randomPoint.X, randomPoint.Y, mouseSpeed[0], 1)
	} else {
		self.MoveTo(randomPoint.X, randomPoint.Y, 3, 1)
	}

	return true
}
func (self *GUI) Init(rect *win.RECT) {
	self.wRect = rect
	self.winRegion = region{int(float64(rect.Left) * dpi), int(float64(rect.Top) * dpi), int(float64(rect.Right-rect.Left) * dpi), int(float64(rect.Bottom-rect.Top) * dpi)}
	if self.picPositionsMap==nil{
		self.picPositionsMap=make(picPosition)
	}


}
func (self *GUI) boundPoints(points []robotgo.Point) []robotgo.Point {
	for i := 0; i < len(points); i++ {
		for {
			if points[i].X >= int(self.wRect.Left) && points[i].X <= int(self.wRect.Right) {
				break
			} else {
				if points[i].X <= int(self.wRect.Left) {
					points[i].X = int(self.wRect.Left) + int(self.wRect.Left) - points[i].X
				}
				if points[i].X >= int(self.wRect.Right) {
					points[i].X = int(self.wRect.Right) - (points[i].X - int(self.wRect.Right))
				}

			}

		}

	}
	for i := 0; i < len(points); i++ {
		for {
			if points[i].Y >= int(self.wRect.Top) && points[i].Y <= int(self.wRect.Bottom) {
				break
			} else {
				if points[i].Y <= int(self.wRect.Top) {
					points[i].Y = int(self.wRect.Top) + int(self.wRect.Top) - points[i].Y
				}
				if points[i].Y >= int(self.wRect.Bottom) {
					points[i].Y = int(self.wRect.Bottom) - (points[i].Y - int(self.wRect.Bottom))
				}

			}

		}

	}
	return points
}

// MoveTo 移动鼠标到某点，传入的参数是该位置在屏幕中的位置而非像素中的位置，也就是要传入没缩放之前的位置
//第三个参数取值1-15，表示鼠标移动速度（不是精确速度）
//第四个参数表示是否对路径进行限定，传入任意int就会限定线路点不超出窗口的范围
func (self *GUI) MoveTo(x1, y1 int, paramList ...int) {
	mouseSpeed := 3
	if len(paramList) > 0 {
		mouseSpeed = paramList[0]
	}
	if mouseSpeed > 15 {
		mouseSpeed = 15
	}
	if mouseSpeed < 1 {
		mouseSpeed = 1
	}

	var method string
	method = methodListAll[rand.Intn(len(methodListAll))]
	//fmt.Println(method)
	//fmt.Println(robotgo.GetMousePos())
	x0, y0 := robotgo.GetMousePos()
	points := GetPointsOnLine(x0, y0, x1, y1, method)
	if len(paramList) > 1 {
		if self.wRect == nil {
			fmt.Println("self.Rect==nil")
			return
		} else {
			points = self.boundPoints(points)
		}
	}
	count := 0

	for i := 0; i < len(points); i++ {
		if count == 0 {
			//实际上sleep的最小间隔大概在15ms左右
			time.Sleep(time.Millisecond * 8)
			//通过修改以下两个数字大概可以改变鼠标的平移速度
			count = rand.Intn(3*mouseSpeed) + 8
			//fmt.Println(count)
		}
		count -= 1
		robotgo.Move(points[i].X, points[i].Y)
	}
}

//截取屏幕，传入的参数是矩形起始位置和长宽
func (self *GUI) winScreenShot(Region region) (gocv.Mat, error) {
	//rect := self.wRect
	//if *rect == (win.RECT{}) {
	//	return gocv.Mat{}, errors.New("wRect=nil,can't get screenShot")
	//}
	bitmap := robotgo.CaptureScreen(Region.Left, Region.Top, Region.Width, Region.Height)
	imageTmp := robotgo.ToImage(bitmap)
	defer robotgo.FreeBitmap(bitmap)
	return cvtImageToMat(imageTmp)
}

//图片匹配，返回最佳匹配区域的左上角以及是否匹配成功
func (self *GUI) matchLocate(source, templ gocv.Mat) (image.Point, bool) {

	var img1, img2 gocv.Mat
	var result, mask gocv.Mat

	result = gocv.NewMatWithSize(source.Rows()-templ.Rows()+1, source.Cols()-templ.Cols()+1, gocv.MatTypeCV32F)
	mask = gocv.Zeros(templ.Rows(), templ.Cols(), gocv.MatTypeCV8U)
	img1 = gocv.NewMat()
	img2 = gocv.NewMat()

	defer func(img2 *gocv.Mat) {
		err := img2.Close()
		if err != nil {
			println(err)
		}
	}(&img2)
	defer func(img1 *gocv.Mat) {
		err := img1.Close()
		if err != nil {
			println(err)
		}
	}(&img1)
	defer func(mask *gocv.Mat) {
		err := mask.Close()
		if err != nil {
			println(err)
		}
	}(&mask)
	defer func(result *gocv.Mat) {
		err := result.Close()
		if err != nil {
			println(err)
		}
	}(&result)

	gocv.CvtColor(source, &img1, gocv.ColorRGBToGray)
	gocv.CvtColor(templ, &img2, gocv.ColorRGBToGray)
	fmt.Println(img1.Size(), img2.Size())
	//tem:=gocv.NewWindow("op")
	gocv.IMWrite("ops.png", img1)
	gocv.IMWrite("opt.png", img2)

	gocv.MatchTemplate(img1, img2, &result, gocv.TmSqdiffNormed, mask)
	gocv.Normalize(result, &result, 0, 1, gocv.NormMinMax)

	minVal, _, minLoc, _ := gocv.MinMaxLoc(result)
	println(minVal, minLoc.X, minLoc.Y)
	if minVal > 0.1 {
		return minLoc, false
	}
	return minLoc, true

}
