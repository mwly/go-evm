package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"sync"

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

func (STP *SpaceTimePyramid) GetPyramid(PiT int) Pyramid {
	Pyr := make([]gocv.Mat, STP.Levels+1)
	for i := range Pyr {
		Pyr[i] = STP.Level[i].SpacialPictures[PiT]
	}
	return Pyr

}

func (STP *SpaceTimePyramid) CreateTimelineAt(row int, col int, level int) {

	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	res := InitATimeline(SpacePics[0].Channels(), len(SpacePics))
	for i, pic := range SpacePics {
		tmp := pic.GetVecdAt(col, row)
		for j := range tmp {
			res[j][i] = tmp[j]
		}
	}
	Pyr.TemporalPictures[row][col] = res
}

type Vecb []float64

func (v Vecb) SetVecbAt(m gocv.Mat, row int, col int) {
	ch := m.Channels()
	for c := 0; c < ch; c++ {
		m.SetDoubleAt(row, col*ch+c, v[c])
	}
}

func (STP *SpaceTimePyramid) ReverseTimelineAt(row int, col int, level int) {
	Pyr := &(STP.Level[level])
	SpacePics := Pyr.SpacialPictures
	for i := range SpacePics {
		Vecb := make(Vecb, 3)
		for j := range Vecb {
			Vecb[j] = Pyr.TemporalPictures[row][col][j][i]
		}
		Vecb.SetVecbAt(Pyr.SpacialPictures[i], row, col)
	}
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

func DoFFTonLine(line *FrequencyLine, i int, r int, c int, ch int, STP *SpaceTimePyramid, WG *sync.WaitGroup) {
	spectrum := fft.FFTReal((*STP).Level[i].TemporalPictures[r][c][ch])
	(*line)[ch] = spectrum
	(*WG).Done()
}

func CreateFrequencyPyrFromSpaceTimePyr(STP SpaceTimePyramid) FrequencyPyramid {
	// Create a 5D Array aka Pyramid
	// level,pictures(rows,cols,chan,point in time)
	level := STP.Levels
	rows := STP.Rows
	cols := STP.Cols
	chanum := STP.Level[0].SpacialPictures[0].Channels()
	frame := len(STP.Level[0].SpacialPictures)
	FP := make([][][]FrequencyLine, level+1)
	var WG sync.WaitGroup
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
					WG.Add(1)
					go DoFFTonLine(&line, i, r, c, ch, &STP, &WG)
					fmt.Printf("start fft for level: %v row: %v col: %v channel: %v \n", i, r, c, ch)
				}
				WG.Wait()
				Fcol[c] = line
			}
			Frows[r] = Fcol
		}
		FP[i] = Frows

	}
	return FrequencyPyramid{level, rows, cols, chanum, frame, FP}
}

func (FP *FrequencyPyramid) CreateSpaceTimePyramidfromFrequencyPyramid(STP SpaceTimePyramid) SpaceTimePyramid {
	MYSTP := STP.Copy()
	var WG sync.WaitGroup
	for i, le := range MYSTP.Level {
		for ro := range le.TemporalPictures {
			for col := range le.TemporalPictures[ro] {
				for ch := range le.TemporalPictures[ro][col] {
					fmt.Printf("start ifft for level: %v row: %v col: %v channel: %v \n", i, ro, col, ch)
					WG.Add(1)
					go func(FP *FrequencyPyramid, i int, ro int, col int, ch int, MYSTP *SpaceTimePyramid, WG *sync.WaitGroup) {
						comp := fft.IFFT((*FP).Pyr[i][ro][col][ch])
						for x := range comp {
							(*MYSTP).Level[i].TemporalPictures[ro][col][ch][x] = cmplx.Abs(comp[x])
						}
						WG.Done()
					}(FP, i, ro, col, ch, &MYSTP, &WG)
				}
			}

		}

	}
	return MYSTP
}

/*
func (FP *FrequencyPyramid) ApplyFilterToFrequencyPyramid(filter Filter) {

}

type Filter struct {
	fsamp   int
	fnyq    int
	numsamp int
	fstart  int
	fend    int
}

func CreateFilter(fsamp int)

func (F *Filter) ApplyToCompl128(pArr *[]complex128) {

}
*/
