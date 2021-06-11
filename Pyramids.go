package main

import (
	"gocv.io/x/gocv"
)

type Pyramid []gocv.Mat

type PyramidLevel struct {
	SpacialPictures []gocv.Mat
}

type TimePyramid struct {
	Levels   int
	Frames   int
	RootRows int
	RootCols int
	Level    []PyramidLevel
}

func CreateTimePyramid(levels int, frames int, rows int, cols int) TimePyramid {
	PyrLevels := make([]PyramidLevel, levels+1)
	for i := range PyrLevels {
		mat := make([]gocv.Mat, frames)
		PyrLevels[i] = PyramidLevel{mat}
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

func (STP *TimePyramid) CreateTimelineFrom(row int, col int, level int) TimeLine {

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

type Vecb []float64

func (v Vecb) SetVecbAt(m gocv.Mat, row int, col int) {
	ch := m.Channels()
	for c := 0; c < ch; c++ {
		m.SetDoubleAt(row, col*ch+c, v[c])
	}
}

func (STP *TimePyramid) InsertTimelineAt(row int, col int, level int, TL *TimeLine) {
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

func (STP *TimePyramid) Copy() TimePyramid {
	return *STP
}

func (STP *TimePyramid) FilterAt(row int, col int, level int, fil Filter, chanum int) {
	TL := STP.CreateTimelineFrom(row, col, level)
	FL := TL.CreateFrequencyLine(chanum)
	fil.ApplyToCompl128(&FL)
	TL = FL.CreateTimeline(chanum)
	STP.InsertTimelineAt(row, col, level, &TL)
}
