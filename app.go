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

func CalcGaussandLapl(img *gocv.Mat, gauss *gocv.Mat, lapl *gocv.Mat) (err error) {
	temp := img.Clone()
	fmt.Println("entered function calcGaussandLapl printing shapes of img, " + fmt.Sprint(img.Size()) + ", gauss: " + fmt.Sprint(gauss.Size()) + ", lapl: " + fmt.Sprint(lapl.Size()))
	gocv.PyrDown(*img, gauss, image.Pt((img.Rows()+1)/2, (img.Cols()+1)/2), gocv.BorderDefault)
	gocv.PyrUp(*gauss, &temp, image.Pt((gauss.Rows()*2), (gauss.Cols()*2)), gocv.BorderDefault)
	fmt.Println("try to subtract temp, " + fmt.Sprint(temp.Size()) + ", from img: " + fmt.Sprint(img.Size()))
	gocv.Subtract(*img, temp, lapl)
	return nil
}

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

	TemporalPyramids := make([]Pyramid, Props.fcount)

	for i := 0; i < Props.fcount; i++ {
		if ok := vid.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}
		if img.Empty() {
			continue
		}

		TemporalPyramids[i] = make([]gocv.Mat, (levels + 1))

		//fmt.Println(img.Size())

		window.IMShow(img)

		fmt.Println(len(TemporalPyramids[i]))
		for j := range TemporalPyramids[i] {

			fmt.Println("Vorgang nr: " + fmt.Sprint(j))

			if j == len(TemporalPyramids[j])-1 {
				fmt.Println("letzter vorgang")
			}

			CalcGaussandLapl(&img, &gaussian_img, &laplacian_img)

			if j == len(TemporalPyramids[i])-1 {
				fmt.Println("entered if")
				laplacian_img = img.Clone()
			} else {
				fmt.Println("entered else")
				img = gaussian_img.Clone()
			}

			TemporalPyramids[i][j] = laplacian_img.Clone()
			defer TemporalPyramids[i][j].Close()
			/*
				gocv.PyrDown(img, &gaussian_img, image.Pt((img.Rows()+1)/2, (img.Cols()+1)/2), gocv.BorderDefault)
				gocv.PyrUp(gaussian_img, &gaussian_img, image.Pt((gaussian_img.Rows()*2), (gaussian_img.Cols()*2)), gocv.BorderDefault)
				gocv.Subtract(img, gaussian_img, &laplacian_img)
			*/
		}

		window2.IMShow(TemporalPyramids[i][0])
		window3.IMShow(TemporalPyramids[i][1])
		window4.IMShow(TemporalPyramids[i][2])
		window5.IMShow(TemporalPyramids[i][3])
		window6.IMShow(TemporalPyramids[i][4])
		//fmt.Println(img.Size())
		window.WaitKey(32)

	}

	for i := range TemporalPyramids {
		for j := range TemporalPyramids[i] {
			fmt.Println(TemporalPyramids[i][j].Size())
		}
	}
	fmt.Println(len(TemporalPyramids))
}
