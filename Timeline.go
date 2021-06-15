package main

import (
	"errors"

	"github.com/mjibson/go-dsp/fft"
)

type TimeLine [][]float64

func InitATimeline(chanum int, len int) TimeLine {
	line := make([][]float64, chanum)
	for i := range line {
		colorSpace := make([]float64, len)
		line[i] = colorSpace
	}
	return line

}

type FrequencyLine [][]complex128

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
	fsamp    int
	fnyq     float64
	fstart   float64
	fend     float64
	Nmin     int
	LevelMin int
}

func CreateFilter(fsamp int, fstart float64, fend float64, LevelMin int) (Filter, error) {
	fnyq := float64(fsamp) / 2
	if fend < 0 && fstart < 0 {
		return Filter{}, errors.New("filter creation: start or end frequenzy isnt positive")
	}
	if fstart > fend {
		return Filter{}, errors.New("filter creation: start frequenzy higher than end frequency")
	}
	if fend > fnyq {
		return Filter{}, errors.New("filter creation: end frequenzy higher than nyquist frequency")
	}
	//I define hereby to oversample with atleast factor 2 therefore the equation for the minimal amount of Samples is
	Nmin := (fsamp * 2)

	return Filter{fsamp, fnyq, fstart, fend, Nmin, LevelMin}, nil
}

func (F *Filter) ApplyToCompl128(pArr *FrequencyLine) {
	//fmt.Printf("Apply Filter to complex signal with length %v \n", len((*pArr)[0]))
	if F.Nmin > len((*pArr)[0]) {
		panic("Number of samples is to small for this Filter")
	}
	df := float64(F.fsamp) / float64(len((*pArr)[0]))
	Nstart := int(F.fstart / df)
	Nend := int(F.fend/df) + 1
	amp := []float64{float64(0), float64(0), float64(80)}
	for n := range (*pArr)[0] {
		if n < Nstart || n > Nend {
			//for ch := range *pArr {
			//(*pArr)[ch][n] = 0 + 0i
			//}
		} else {
			for ch := range *pArr {
				(*pArr)[ch][n] += complex(amp[ch]*real((*pArr)[ch][n]), amp[ch]*imag((*pArr)[ch][n]))
			}
		}
	}

}

type RoomInPyr struct {
	row    int
	col    int
	level  int
	fil    Filter
	chanum int
}
