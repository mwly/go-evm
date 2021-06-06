package main

import (
	"fmt"

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

	window2 := gocv.NewWindow("2")
	defer window2.Close()

	window3 := gocv.NewWindow("3")
	defer window2.Close()

	window4 := gocv.NewWindow("4")
	defer window2.Close()

	window5 := gocv.NewWindow("5")
	defer window2.Close()

	window6 := gocv.NewWindow("6")
	defer window2.Close()

	//window7 := gocv.NewWindow("7")
	//defer window2.Close()

	img := gocv.NewMat()
	defer img.Close()

	gaussian_img := gocv.NewMat()
	defer gaussian_img.Close()

	laplacian_img := gocv.NewMat()
	defer laplacian_img.Close()

	//imges := []gocv.Mat
	Props := InitVideoProperties(*vid)

	fmt.Println(img.Channels())
	fmt.Println(Props)

	TemporalPyramid := make([]Pyramid, Props.fcount)

	CroppingRect := GetCroppingRect(Props)

	for i := 0; i < Props.fcount; i++ {
		if ok := vid.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}
		if img.Empty() {
			continue
		}

		img_temp := img.Region(CroppingRect)
		img = img_temp.Clone()
		window.IMShow(img)

		CreatePyramid(&TemporalPyramid[i], img, levels)

		img = ReconstructImageFromPyramid(TemporalPyramid[i]).Clone()

		window2.IMShow(TemporalPyramid[i][0])
		window3.IMShow(TemporalPyramid[i][1])
		window4.IMShow(TemporalPyramid[i][2])
		window5.IMShow(TemporalPyramid[i][3])
		window6.IMShow(img)
		//fmt.Println(img.Size())
		//fmt.Println(img.Size())
		window.WaitKey(32)

	}
	/*
		for i := range TemporalPyramids {
			for j := range TemporalPyramids[i] {
				fmt.Println(TemporalPyramids[i][j].Size())
			}
		}
	*/

	//fft.FFTReal()
	fmt.Println(len(TemporalPyramid))
	window.WaitKey(-1)
}
