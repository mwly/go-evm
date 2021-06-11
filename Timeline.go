package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/mjibson/go-dsp/fft"
)

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
					fmt.Printf("\r start fft for level: %v row: %v col: %v channel: %v ", i, r, c, ch)
				}
				WG.Wait()
				Fcol[c] = line
			}
			Frows[r] = Fcol
		}
		FP[i] = Frows

	}
	fmt.Printf("\n")
	return FrequencyPyramid{level, rows, cols, chanum, frame, FP}
}

func ReconComplexSignal(idft []complex128) []float64 {
	data := make([]float64, len(idft))
	for i := 0; i < len(data); i++ {
		data[i] = real(idft[i])
	}
	return data
}

func (FP *FrequencyPyramid) CreateSpaceTimePyramidfromFrequencyPyramid(STP SpaceTimePyramid) SpaceTimePyramid {
	MYSTP := STP.Copy()
	var WG sync.WaitGroup
	for i, le := range MYSTP.Level {
		for ro := range le.TemporalPictures {
			for col := range le.TemporalPictures[ro] {
				for ch := range le.TemporalPictures[ro][col] {
					fmt.Printf("\r start ifft for level: %v row: %v col: %v channel: %v ", i, ro, col, ch)
					WG.Add(1)
					go func(FP *FrequencyPyramid, i int, ro int, col int, ch int, MYSTP *SpaceTimePyramid, WG *sync.WaitGroup) {

						comp := fft.IFFT((*FP).Pyr[i][ro][col][ch])

						(*MYSTP).Level[i].TemporalPictures[ro][col][ch] = ReconComplexSignal(comp)

						WG.Done()
					}(FP, i, ro, col, ch, &MYSTP, &WG)
				}

			}

		}

	}
	WG.Wait()
	fmt.Printf("\n")
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
