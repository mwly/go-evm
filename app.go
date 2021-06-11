package main

import (
	"fmt"
	"sync"

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

var file string = "test-face.mp4"
var levels int = 4

func main() {

	vid, err := gocv.VideoCaptureFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer vid.Close()

	Props := InitVideoProperties(*vid)

	Egypt := make([]Pyramid, Props.fcount)

	CroppingRect := GetCroppingRect(Props)

	SpaTiPyr := CreateTimePyramid(levels, Props.fcount, CroppingRect.Dx(), CroppingRect.Dy())

	OutPut, err := gocv.VideoWriterFile(("processed_" + file), vid.CodecString(), vid.Get(gocv.VideoCaptureFPS), CroppingRect.Dx(), CroppingRect.Dy(), true)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer OutPut.Close()

	OutputFrame := gocv.NewMat()
	defer OutputFrame.Close()

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

	window_result := gocv.NewWindow("Result")
	defer window2.Close()

	img := gocv.NewMat()
	defer img.Close()

	gaussian_img := gocv.NewMat()
	defer gaussian_img.Close()

	laplacian_img := gocv.NewMat()
	defer laplacian_img.Close()

	//imges := []gocv.Mat

	fmt.Println(img.Channels())
	fmt.Println(Props)

	for i := 0; i < Props.fcount; i++ {
		if ok := vid.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", file)
			return
		}
		if img.Empty() {
			continue
		}

		img_temp := img.Region(CroppingRect)
		img = ImageTo64float(img_temp)

		Egypt[i] = CreatePyramid(img, levels)

		SpaTiPyr.AddPyramid(Egypt[i], i)

		OutputFrame = ReconstructImageFromPyramid(Egypt[i]).Clone()

		OutPut.Write(ImageTo8Int(OutputFrame))

		//Result.ConvertTo(&Result, gocv.MatTypeCV8UC3)

		//fft.FFT2Real()

		window2.IMShow(ImageTo8Int(Egypt[i][0].Clone()))
		window3.IMShow(ImageTo8Int(Egypt[i][1].Clone()))
		window4.IMShow(ImageTo8Int(Egypt[i][2].Clone()))
		window5.IMShow(ImageTo8Int(Egypt[i][3].Clone()))
		window6.IMShow(ImageTo8Int(Egypt[i][4].Clone()))
		//window7.IMShow(TemporalPyramid[i][5])
		window_result.IMShow(ImageTo8Int(OutputFrame))
		window.IMShow(ImageTo8Int(img.Clone()))
		//fmt.Println(img.Size())
		//fmt.Println(img.Size())
		window.WaitKey(32)

	}

	fil := Filter{}

	newTiPyr := SpaTiPyr.Copy()

	var WG sync.WaitGroup

	for level, levels := range newTiPyr.Level {
		for row := 0; row < levels.SpacialPictures[0].Rows(); row++ {
			for col := 0; col < levels.SpacialPictures[0].Cols(); col++ {
				WG.Add(1)
				go newTiPyr.FilterAt(row, col, level, fil, levels.SpacialPictures[0].Channels(), &WG)
			}
		}
	}
	WG.Wait()

	fmt.Println("STP created printing to movie")

	FFTOutPut, err := gocv.VideoWriterFile(("FFT_processed_" + file), vid.CodecString(), vid.Get(gocv.VideoCaptureFPS), CroppingRect.Dx(), CroppingRect.Dy(), true)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer FFTOutPut.Close()
	for PiT := 0; PiT < newTiPyr.Frames; PiT++ {
		FFTOutPut.Write(ImageTo8Int(*ReconstructImageFromPyramid(newTiPyr.GetPyramid(PiT))))
	}

	fmt.Println(len(Egypt))
	window.WaitKey(-1)
	fmt.Println("lel")
}
