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

var file string = "subway.mp4"
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

	SpaTiPyr := CreateTimePyramid(levels, Props.fcount, CroppingRect.Dy(), CroppingRect.Dx(), 3)

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

		BGR := gocv.Split(OutputFrame)

		gocv.Merge(BGR, &OutputFrame)

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

	numworkers := 5

	ch := make(chan RoomInPyr, (2)*newTiPyr.RootRows*newTiPyr.RootCols)
	fmt.Printf("made channel with len %v \n", 2*newTiPyr.RootRows*newTiPyr.RootCols)

	newTiPyr.RGB2Gray()

	newTiPyr.PrintShape()

	for z := 0; z < numworkers; z++ {
		WG.Add(1)

		fmt.Printf("Create worker number: %v \n", z)
		go func(WG *sync.WaitGroup) {
			for Room := range ch {
				newTiPyr.FilterGrayAt(Room.row, Room.col, Room.level, Room.fil, Room.chanum)
			}
			WG.Done()
		}(&WG)
	}
	WG.Add(1)
	go func(WG *sync.WaitGroup) {
		i := 0

		for Room := range ch {
			fmt.Printf("\rworker 11 did his %v task", i)
			newTiPyr.FilterGrayAt(Room.row, Room.col, Room.level, Room.fil, Room.chanum)
			i += 1
		}
		fmt.Printf("\n")
		WG.Done()
	}(&WG)

	for level, levels := range newTiPyr.Level {

		for row := 0; row < levels.SpacialPictures[0].Rows(); row++ {
			for col := 0; col < levels.SpacialPictures[0].Cols(); col++ {
				ch <- RoomInPyr{row, col, level, fil, 3}
			}
		}
	}
	close(ch)
	WG.Wait()
	newTiPyr.Gray2RGB()
	fmt.Println("STP created printing to movie")
	newTiPyr.PrintShape()

	FFTOutPut, err := gocv.VideoWriterFile(("FFT_processed_" + file), vid.CodecString(), vid.Get(gocv.VideoCaptureFPS), CroppingRect.Dx(), CroppingRect.Dy(), true)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer FFTOutPut.Close()
	for PiT := 0; PiT < newTiPyr.Frames; PiT++ {
		Pyr := newTiPyr.GetPyramid(PiT)
		RecImg := ReconstructImageFromPyramid(Pyr)
		IntImg := ImageTo8Int(*RecImg)
		//gocv.CvtColor(IntImg, &IntImg, gocv.ColorGrayToBGR)

		FFTOutPut.Write((IntImg))

		BGR := gocv.Split(IntImg)

		window.IMShow(IntImg)
		window2.IMShow(BGR[0])
		window3.IMShow(BGR[1])
		window4.IMShow(BGR[2])
		window.WaitKey(32)

		//FFTOutPut.Write(ImageTo8Int(*ReconstructImageFromPyramid(newTiPyr.GetPyramid(PiT))))
	}
	window.WaitKey(-1)
	fmt.Println(len(Egypt))
	fmt.Println("lel")
}
