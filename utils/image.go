package utils

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
)

func cvtImageToMat(img image.Image) (gocv.Mat, error) {

	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}
func imageToGray(img image.Image,config image.Config) [][]uint8 {
	var res  [][]uint8
	var a int8=1
	b:=a<<8
	fmt.Println(b)
	return res
}
