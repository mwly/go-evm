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
var levels int = 6

func GetCroppingRect(Props VideoProperties) image.Rectangle {
	width := 2
	height := 2
	for width < (Props.width)/2 {
		width = width * 2
	}
	for height < (Props.height)/2 {
		height = height * 2
	}
	hdc := (Props.height - height) / 2
	wdc := (Props.width - width) / 2

	return image.Rect(wdc, hdc, wdc+width, hdc+height)
}

func MakePyrKernel(img gocv.Mat, up bool) image.Point {
	Rows := img.Rows()
	Cols := img.Cols()

	if up {
		if Rows%2 == 1 {
			Rows += 1
		}
		if Cols%2 == 1 {
			Cols += 1
		}
	}

	return image.Pt(Rows/2, Cols/2)

}

func CalcGaussandLapl(img *gocv.Mat, gauss *gocv.Mat, lapl *gocv.Mat) (err error) {
	temp := img.Clone()

	//fmt.Println("entered function calcGaussandLapl printing shapes of img, " + fmt.Sprint(img.Size()) + ", gauss: " + fmt.Sprint(gauss.Size()) + ", lapl: " + fmt.Sprint(lapl.Size()))

	//   (img.Rows())/2, (img.Cols())/2

	gocv.PyrDown(*img, gauss, MakePyrKernel(*img, false), gocv.BorderDefault)
	gocv.PyrUp(*gauss, &temp, image.Pt(gauss.Rows()*2, gauss.Cols()*2), gocv.BorderDefault)

	//fmt.Println("went Pyr up and down shapes of img, " + fmt.Sprint(img.Size()) + ", gauss: " + fmt.Sprint(gauss.Size()) + ", lapl: " + fmt.Sprint(lapl.Size()) + ", temp: " + fmt.Sprint(temp.Size()))
	//fmt.Println("try to subtract temp, " + fmt.Sprint(temp.Size()) + ", from img: " + fmt.Sprint(img.Size()))
	gocv.Subtract(*img, temp, lapl)
	return nil
}

func CreatePyramid(Pyr *Pyramid, img gocv.Mat, levels int) {

	*Pyr = make([]gocv.Mat, (levels + 1))

	gaussian_img := gocv.NewMat()
	defer gaussian_img.Close()

	laplacian_img := gocv.NewMat()
	defer laplacian_img.Close()

	//fmt.Println(img.Size())
	//fmt.Println(len(*Pyr))

	for j := range *Pyr {

		//fmt.Println("Vorgang nr: " + fmt.Sprint(j))

		CalcGaussandLapl(&img, &gaussian_img, &laplacian_img)

		if j == len(*Pyr)-1 {
			//fmt.Println("entered if")
			laplacian_img = img.Clone()
		} else {
			//fmt.Println("entered else")
			img = gaussian_img.Clone()
		}

		(*Pyr)[j] = laplacian_img.Clone()
	}
}

func ReconstructImageFromPyramid(Pyr Pyramid) *gocv.Mat {
	Result := Pyr[len(Pyr)-1].Clone()

	for i := len(Pyr) - 2; i >= 0; i-- {
		gocv.PyrUp(Result, &Result, image.Pt(Result.Rows()*2, Result.Cols()*2), gocv.BorderDefault)
		gocv.Add(Result, Pyr[i], &Result)
		//fmt.Println("Added Pyr with shape (" + fmt.Sprint(Result.Size()) + ") and Result with shape (" + fmt.Sprint(Result.Size()) + ")")
	}
	return &Result
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

		CreatePyramid(&TemporalPyramids[i], img, levels)

		img = ReconstructImageFromPyramid(TemporalPyramids[i]).Clone()

		window2.IMShow(TemporalPyramids[i][0])
		window3.IMShow(TemporalPyramids[i][1])
		window4.IMShow(TemporalPyramids[i][2])
		window5.IMShow(TemporalPyramids[i][3])
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
	fmt.Println(len(TemporalPyramids))
	window.WaitKey(-1)
}
