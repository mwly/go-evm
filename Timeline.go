package main

import (
	"github.com/mjibson/go-dsp/fft"
)

type TimeLine [][]float64

func InitATimeline(chanum int, len int) TimeLine {
	line := make([][]float64, chanum)
	colorSpace := make([]float64, len)
	for i := range line {
		line[i] = colorSpace
	}
	return line

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

func (TL *TimeLine) CreateFrequencyLine(chanum int) FrequencyLine {

	line := make(FrequencyLine, chanum)
	for ch := range line {
		//iterate along the channels of the column

		DoFFTonLine(&line, ch, TL)
	}
	return line
}

func DoFFTonLine(line *FrequencyLine, ch int, TL *TimeLine) {
	spectrum := fft.FFTReal((*TL)[ch])
	(*line)[ch] = spectrum
}

func (FL *FrequencyLine) CreateTimeline(chanum int) TimeLine {
	TL := make(TimeLine, chanum)
	for ch := range *FL {
		func(FL *FrequencyLine, ch int, TL *TimeLine) {

			comp := fft.IFFT((*FL)[ch])

			(*TL)[ch] = ReconComplexSignal(comp)
		}(FL, ch, &TL)
	}
	return TL
}

func ReconComplexSignal(idft []complex128) []float64 {
	data := make([]float64, len(idft))
	for i := 0; i < len(data); i++ {
		data[i] = real(idft[i])
	}
	return data
}

type Filter struct {
	fsamp   int
	fnyq    int
	numsamp int
	fstart  int
	fend    int
}

func CreateFilter(fsamp int) {

}

func (F *Filter) ApplyToCompl128(pArr *FrequencyLine) {

}

type RoomInPyr struct {
	row    int
	col    int
	level  int
	fil    Filter
	chanum int
}
