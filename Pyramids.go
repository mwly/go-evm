package main

import "gocv.io/x/gocv"

type Pyramid []gocv.Mat

type TimeLine [][]float64

func InitATimeline(chanum int, len int) TimeLine {
	tmp := make([]float64, len)
	line := make([][]float64, chanum)
	for i := range line {
		line[i] = tmp
	}
	return line

}

func CreateTimeline(TP []Pyramid, row int, col int, level int) TimeLine {
	res := InitATimeline(3, len(TP))

	for i, Pyr := range TP {

		tmp := Pyr[level].GetVecdAt(row, col)
		for j := range tmp {
			res[j][i] = tmp[j]
		}
	}
	return res
}

type PyramidLevel struct {
	SpacialPictures  []gocv.Mat
	TemporalPictures [][]TimeLine
}

type SpaceTimePyramid struct {
	Levels int
	Frames int
	Level  []PyramidLevel
}
