package evm

import (
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

	gocv.PyrDown(*img, gauss, MakePyrKernel(*img, false), gocv.BorderDefault)
	gocv.PyrUp(*gauss, &temp, image.Pt(gauss.Rows()*2, gauss.Cols()*2), gocv.BorderDefault)

	gocv.Subtract(*img, temp, lapl)
	return nil
}

func ImageTo64float(img gocv.Mat) gocv.Mat {
	img.ConvertTo(&img, gocv.MatTypeCV64FC3)
	return img
}

func ImageTo8Int(img gocv.Mat) gocv.Mat {
	img.ConvertTo(&img, gocv.MatTypeCV8UC3)
	return img
}

func ReconstructImageFromPyramid(Pyr Pyramid) *gocv.Mat {
	Result := Pyr[len(Pyr)-1].Clone()

	for i := len(Pyr) - 2; i >= 0; i-- {
		gocv.PyrUp(Result, &Result, image.Pt(Result.Rows()*2, Result.Cols()*2), gocv.BorderDefault)
		//fmt.Println("Try to add Pyr with shape (" + fmt.Sprint(Pyr[i].Size()) + ") and Result with shape (" + fmt.Sprint(Result.Size()) + ")")
		//fmt.Println("Try to add Pyr with type (" + fmt.Sprint(Pyr[i].Type().String()) + ") and Result with shape (" + fmt.Sprint(Result.Type().String()) + ")")
		//fmt.Println("Try to add Pyr with channels (" + fmt.Sprint(Pyr[i].Channels()) + ") and Result with shape (" + fmt.Sprint(Result.Channels()) + ")")
		gocv.Add(Result, Pyr[i], &Result)
	}
	return &Result
}
