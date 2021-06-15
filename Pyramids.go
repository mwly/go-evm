package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Pyramid []gocv.Mat

type PyramidLevel struct {
	SpacialPictures     []gocv.Mat
	SpacialPicturesGray [][]gocv.Mat
}

type TimePyramid struct {
	Levels   int
	Frames   int
	RootRows int
	RootCols int
	Level    []PyramidLevel
}

func CreateTimePyramid(levels int, frames int, rows int, cols int, chanum int) TimePyramid {
	PyrLevels := make([]PyramidLevel, levels+1)
	for i := range PyrLevels {
		mat := make([]gocv.Mat, frames)
		graymat := make([][]gocv.Mat, chanum)
		for x := range graymat {
			a_mat := make([]gocv.Mat, frames)
			graymat[x] = a_mat
		}
		PyrLevels[i] = PyramidLevel{mat, graymat}
	}
	return TimePyramid{levels, frames, rows, cols, PyrLevels}

}

func (STP *TimePyramid) AddPyramid(Pyr Pyramid, PiT int) {
	//PiT point int Time
	for i, p := range Pyr {
		STP.Level[i].SpacialPictures[PiT] = p
	}
}

func (STP *TimePyramid) GetPyramid(PiT int) Pyramid {
	Pyr := make([]gocv.Mat, STP.Levels+1)
	for i := range Pyr {
		Pyr[i] = STP.Level[i].SpacialPictures[PiT]
	}
	return Pyr

}

func (STP *TimePyramid) CreateTimelineFromRGB(row int, col int, level int) TimeLine {

	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	resultingTimeline := InitATimeline(SpacePics[0].Channels(), len(SpacePics))
	for i, pic := range SpacePics {
		tmp := pic.GetVecdAt(row, col)
		for j := range tmp {
			resultingTimeline[j][i] = tmp[j]
		}
	}
	return resultingTimeline
}

func (STP *TimePyramid) CreateTimelineFromGray(row int, col int, level int, chanum int) TimeLine {

	Pyr := &(STP.Level[level])
	GrayPics := Pyr.SpacialPicturesGray
	resultingTimeline := InitATimeline(chanum, len(GrayPics[0]))
	for i := range GrayPics[0] {
		for j := 0; j < chanum; j++ {
			resultingTimeline[j][i] = GrayPics[j][i].GetDoubleAt(row, col)
		}
	}
	return resultingTimeline
}

type Vecb []float64

func (v Vecb) SetVecbAt(m gocv.Mat, row int, col int) {
	ch := m.Channels()
	for c := 0; c < ch; c++ {
		m.SetDoubleAt(row, col*ch+c, v[c])
	}
}

func (STP *TimePyramid) InsertTimelineAtRGB(row int, col int, level int, TL *TimeLine) {
	//fmt.Printf("inserting Pixels at row: %v col%v and level %v ", row, col, level)
	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	for i := range SpacePics {
		Vecb := make(Vecb, 3)
		for j := range Vecb {
			Vecb[j] = (*TL)[j][i]
		}
		Vecb.SetVecbAt(Pyr.SpacialPictures[i], row, col)
	}
	//fmt.Printf(": Done\n")
}

func (STP *TimePyramid) InsertTimelineAtGray(row int, col int, level int, TL *TimeLine) {
	//fmt.Printf("inserting Pixels at row: %v col%v and level %v ", row, col, level)
	Pyr := &(STP.Level[level])
	GrayPics := Pyr.SpacialPicturesGray
	for i := range GrayPics[0] {
		for j := range GrayPics {
			//fmt.Printf("\n(length, depth) of Timeline l: %v d: %v \n", len(*TL), len((*TL)[j]))
			//fmt.Printf("\n(length, depth) of gray l: %v d: %v \n", len(Pyr.SpacialPicturesGray), len(Pyr.SpacialPicturesGray[j]))
			Pyr.SpacialPicturesGray[j][i].SetDoubleAt(row, col, (*TL)[j][i])
		}
	}
	//fmt.Printf(": Done\n")
}

func (STP *TimePyramid) Copy() TimePyramid {
	return *STP
}

func (STP *TimePyramid) FilterRGBAt(row int, col int, level int, fil Filter, chanum int) {
	TL := STP.CreateTimelineFromRGB(row, col, level)
	FL := TL.CreateFrequencyLine(chanum)
	fil.ApplyToCompl128(&FL)
	TL = FL.CreateTimeline(chanum)
	STP.InsertTimelineAtRGB(row, col, level, &TL)
}

func (STP *TimePyramid) FilterGrayAt(row int, col int, level int, fil Filter, chanum int) {
	TL := STP.CreateTimelineFromGray(row, col, level, chanum)
	FL := TL.CreateFrequencyLine(chanum)
	fil.ApplyToCompl128(&FL)
	TL = FL.CreateTimeline(chanum)
	STP.InsertTimelineAtGray(row, col, level, &TL)
}

func (STP *TimePyramid) RGB2Gray() {
	for i := range STP.Level {
		LEVEL := &STP.Level[i]
		for j := range LEVEL.SpacialPictures {
			BGR := gocv.Split(LEVEL.SpacialPictures[j])
			for k := range BGR {
				LEVEL.SpacialPicturesGray[k][j] = BGR[k]
			}
		}
	}
}

func (STP *TimePyramid) Gray2RGB() {
	for i := range STP.Level {
		LEVEL := &STP.Level[i]
		for j := range LEVEL.SpacialPictures {
			BGR := make([]gocv.Mat, 3)
			for k := range BGR {
				BGR[k] = LEVEL.SpacialPicturesGray[k][j]
				//gocv.CvtColor(BGR[k], &(BGR[k]), gocv.ColorBGRToGray)
			}
			//xor := BGR[0].Clone()
			//gocv.BitwiseXor(BGR[0], BGR[1], &xor)
			//if gocv.CountNonZero(xor) == 0 {
			//	panic("error")
			//}
			//fmt.Printf("Tested if we are about to merge the same pictures at level %v and frame %v\n", i, j)
			gocv.Merge(BGR, &LEVEL.SpacialPictures[j])
			//fmt.Printf("Merged BGR with len : %v into image with %v channels \n", len(BGR), LEVEL.SpacialPictures[j].Channels())
		}
	}

}

func (STP *TimePyramid) PrintShape() {
	fmt.Printf("Printing shape of Pyramids:\n")
	for i := range STP.Level {
		fmt.Printf("At Levle %v the RGBshape is %v and the Gray shape is %v\n", i, STP.Level[i].SpacialPictures[0].Size(), STP.Level[i].SpacialPicturesGray[0][0].Size())
	}

}
