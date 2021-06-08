package main

import (
	"math"

	"gocv.io/x/gocv"
)

type Pyramid []gocv.Mat

type TimeLine [][]float64

func InitATimeline(chanum int, len int) TimeLine {
	line := make([][]float64, chanum)
	colorSpace := make([]float64, len)
	for i := range line {
		line[i] = colorSpace
	}
	return line

}

func (STP *SpaceTimePyramid) CreateTimelineAt(row int, col int, level int) {

	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	res := InitATimeline(3, len(SpacePics))
	for i, pic := range SpacePics {
		tmp := pic.GetVecdAt(row, col)
		for j := range tmp {
			res[j][i] = tmp[j]
		}
	}
	Pyr.TemporalPictures[row][col] = res
}

func (STP *SpaceTimePyramid) GetTimelineAt(row int, col int, level int) TimeLine {
	return STP.Level[level].TemporalPictures[row][col]
}

type PyramidLevel struct {
	SpacialPictures  []gocv.Mat
	TemporalPictures [][]TimeLine
}

type SpaceTimePyramid struct {
	Levels int
	Frames int
	Rows   int
	Cols   int
	Level  []PyramidLevel
}

func CreateSpaceTimePyramid(levels int, frames int, rows int, cols int) SpaceTimePyramid {
	PyrLevels := make([]PyramidLevel, levels+1)
	for i := range PyrLevels {
		mat := make([]gocv.Mat, frames)
		TPics := make([][]TimeLine, rows/(int(math.Pow(float64(2), float64(i)))))
		TLines := make([]TimeLine, cols/(int(math.Pow(float64(2), float64(i)))))
		for j := range TPics {
			TPics[j] = TLines
		}

		PyrLevels[i] = PyramidLevel{mat, TPics}
	}

	return SpaceTimePyramid{levels, frames, rows, cols, PyrLevels}

}

func (STP *SpaceTimePyramid) AddPyramid(Pyr Pyramid, PiT int) {
	//PiT point int Time
	for i, p := range Pyr {
		STP.Level[i].SpacialPictures[PiT] = p
	}
}
