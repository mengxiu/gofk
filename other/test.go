package other

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"
)

func main() {
img1:=readPng("test.png")
img2:=readPng("test1.png")
t1:=time.Now()
match(img1,img2,0.5)
fmt.Println(time.Since(t1))
}
func readPng(str string	)image.Image  {
	f,err:=os.Open(str)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	img,err2:=png.Decode(f)
	if err2 != nil {
		fmt.Println(err2)
		return nil
	}
	return img

}
func match(img1, img2 image.Image, p float64) (int, int, bool) {
	var wg sync.WaitGroup
	ori := ToMat(img1)
	templ := ToMat(img2)
	oriDiffX, oriDiffY := ToDiffMat(ori)
	templDiffX, templDiffY := ToDiffMat(templ)

	oridX, oridY := (ori).Size()
	templdX, templdY := (templ).Size()
	resX := make([][]float64, oridY-templdY)
	for i := range resX {
		resX[i] = make([]float64, oridX-templdX)
	}
	resY := make([][]float64, oridY-templdY)
	for i := range resY {
		resY[i] = make([]float64, oridX-templdX)
	}
	for i := 0; i < oridY-templdY; i++ {
		for j := 0; j < oridX-templdX; j++ {
			wg.Add(1)
			go matchMethod(i, j, templdX,templdY,&oriDiffX, &templDiffX,&resX,&wg)
			wg.Add(1)
			go matchMethod(i, j, templdX,templdY,&oriDiffY, &templDiffY,&resY,&wg)


		}
	}
	wg.Wait()
	x:=0.0
	y:=0.0
	for i := 0; i < oridY-templdY; i++ {
		for j := 0; j < oridX-templdX; j++ {
			if resX[i][j]>x{
				x=resX[i][j]
				println(i,j)
			}
			if resY[i][j]>y{
				y=resY[i][j]
				println(i,j)
			}
		}
	}
	println(x,y)
	return 0, 0, false
}
func matchMethod(i, j ,dX,dY int, ori, templ *imgMat,res *[][]float64, wg *sync.WaitGroup) {
	count:=0.0
	for m:=0;m<dY;m++{
		for n:=0;n<dX;n++{
			if (*ori)[i+m][j+n]==(*templ)[m][n]{
				count+=1
			}

		}
	}
	(*res)[i][j]=count/float64(dX*dY)
	wg.Done()
}

func ToGray(img image.Image) image.Image {
	fmt.Println(img.Bounds())
	var gray *image.Gray
	gray = image.NewGray(image.Rectangle{
		Min: image.Point{img.Bounds().Min.X, img.Bounds().Min.Y},
		Max: image.Point{img.Bounds().Max.X, img.Bounds().Max.Y},
	})
	for i := (img.Bounds().Min).X; i < (img.Bounds().Max).X; i++ {

		for j := (img.Bounds().Min).Y; j < (img.Bounds().Max).Y; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			grayPix := uint8((r*76 + g*150 + b*30) >> 16)
			//grayPix:=uint8(0)
			gray.Set(i, j, color.Color(color.RGBA{grayPix, grayPix, grayPix, (uint8(255))}))

		}

	}
	return gray
}

type imgMat [][]uint8

func (self *imgMat) Size() (int, int) {
	dy := len(*self)
	if dy > 0 {
		return len((*self)[0]), dy
	} else {
		return 0, dy
	}
}
func ToDiffMat(ori imgMat) (imgMat, imgMat) {
	dx, dy := ori.Size()
	matX := make([][]uint8, dy)
	for j := range matX {
		matX[j] = make([]uint8, dx)
	}
	matY := make([][]uint8, dy)
	for j := range matY {
		matY[j] = make([]uint8, dx)
	}
	for i := 0; i < dy; i++ {
		for j := 0; j < dx-1; j++ {
			if ori[i][j] < ori[i][j+1] {
				matX[i][j] = 1
			} else {
				matX[i][j] = 0
			}

		}
		matX[i][dx-1] = 0
	}
	for i := 0; i < dx; i++ {
		for j := 0; j < dy-1; j++ {
			if ori[j][i] < ori[j+1][i] {
				matY[j][i] = 1
			} else {
				matY[j][i] = 0
			}

		}
		matY[dy-1][i] = 0
	}

	return matX, matY
}
func ToMat(img image.Image) imgMat {
	bounds := img.Bounds()
	col, row := bounds.Dx(), bounds.Dy()
	var mat imgMat
	mat = make([][]uint8, row)
	for j := range mat {
		mat[j] = make([]uint8, col)
	}
	for i := 0; i < col; i++ {
		for j := 0; j < row; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			mat[j][i] = uint8((r*76 + g*150 + b*30) >> 16)
		}

	}
	return mat

}
func SaveToPng(filename string, img image.Image) error {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}
