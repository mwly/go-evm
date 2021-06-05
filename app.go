package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

type VideoProperties struct {
	fps      int
	height   int
	width    int
	fcount   int
	channels int
	isRGB    int
}

func InitVideoProperties(vid gocv.VideoCapture) VideoProperties {
	fps := vid.Get(gocv.VideoCaptureFPS)
	height := vid.Get(gocv.VideoCaptureFrameHeight)
	width := vid.Get(gocv.VideoCaptureFrameWidth)
	fcount := vid.Get(gocv.VideoCaptureFrameCount)
	channels := vid.Get(gocv.VideoCaptureChannel)
	isRGB := vid.Get(gocv.VideoCaptureConvertRGB)

	Props := VideoProperties{fps: int(fps), height: int(height), width: int(width), fcount: int(fcount), channels: int(channels), isRGB: int(isRGB)}

	return Props
}

type Pyramid []gocv.Mat

var file string = "./test-face.mp4"
var levels int = 4

func main() {

	vid, err := gocv.VideoCaptureFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer vid.Close()

	window := gocv.NewWindow("Std_out")
	defer window.Close()

	window2 := gocv.NewWindow("Debug")
	defer window2.Close()

	img := gocv.NewMat()
	defer img.Close()

	img2 := gocv.NewMat()
	defer img2.Close()

	//imges := []gocv.Mat
	Props := InitVideoProperties(*vid)

	fmt.Println(img.Channels())
	fmt.Println(Props)

	Pyramids := make([]Pyramid, Props.fcount)

	for i := 0; i < Props.fcount; i++ {
		if ok := vid.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}
		if img.Empty() {
			continue
		}

		Pyramids[i] = make([]gocv.Mat, levels)

		//fmt.Println(img.Size())

		Pyramids[i][0] = img

		window.IMShow(img)

		for j := range Pyramids[i] {
			Pyramids[i][j] = img

			gocv.PyrDown(img, &img, image.Pt((img.Rows()+1)/2, (img.Cols()+1)/2), gocv.BorderDefault)
		}

		window2.IMShow(img)
		//fmt.Println(img.Size())
		window.WaitKey(32)
	}

}
