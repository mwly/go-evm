package main

import (
	"math"

	"github.com/mjibson/go-dsp/fft"
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

func (STP *SpaceTimePyramid) CreateTimelineAt(row int, col int, level int) {

	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	res := InitATimeline(SpacePics[0].Channels(), len(SpacePics))
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

func (STP *SpaceTimePyramid) Copy() SpaceTimePyramid {
	return *STP
}

type FrequencyLine [][]complex128

type FrequencyPyramid struct {
	Levels   int
	Rootrows int
	Rootcols int
	Chanum   int
	Frame    int
	Pyr      [][][]FrequencyLine
}

func CreateFrequencyPyrFromSpaceTimePyr(STP SpaceTimePyramid) FrequencyPyramid {
	// Create a 5D Array aka Pyramid
	// level,pictures(rows,cols,chan,point in time)
	level := STP.Levels
	rows := STP.Rows
	cols := STP.Cols
	chanum := STP.Level[0].SpacialPictures[0].Channels()
	frame := len(STP.Level[0].SpacialPictures)
	FP := make([][][]FrequencyLine, level)
	for i := range FP {
		// iterate across the levels of the pyramid
		thisrow := rows / int(math.Pow(float64(2), float64(i)))
		thiscol := cols / int(math.Pow(float64(2), float64(i)))
		Frows := make([][]FrequencyLine, thisrow)
		for r := range Frows {
			//iterate along the rows of the pictures
			Fcol := make([]FrequencyLine, thiscol)
			for c := range Fcol {
				//iterate alon the colums of the rows
				line := make(FrequencyLine, chanum)
				for ch := range line {
					//iterate along the channels of the column
					spectrum := fft.FFTReal(STP.Level[i].TemporalPictures[r][c][ch])
					line[ch] = spectrum
				}
				Fcol[c] = line
			}
			Frows[r] = Fcol
		}
		FP[i] = Frows

	}

	return FrequencyPyramid{level, rows, cols, chanum, frame, FP}

}
