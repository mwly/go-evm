package main

import (
	"image"

	"gocv.io/x/gocv"
)

func GetCroppingRect(Props VideoProperties) image.Rectangle {
	width := 2
	height := 2
	for width < (Props.width)/2 {
		width = width * 2
	}
	for height < (Props.height)/2 {
		height = height * 2
	}
	hdc := (Props.height - height) / 2
	wdc := (Props.width - width) / 2

	return image.Rect(wdc, hdc, wdc+width, hdc+height)
}

func MakePyrKernel(img gocv.Mat, up bool) image.Point {
	Rows := img.Rows()
	Cols := img.Cols()

	if up {
		if Rows%2 == 1 {
			Rows += 1
		}
		if Cols%2 == 1 {
			Cols += 1
		}
	}

	return image.Pt(Rows/2, Cols/2)

}

func CalcGaussandLapl(img *gocv.Mat, gauss *gocv.Mat, lapl *gocv.Mat) (err error) {
	temp := img.Clone()

	//fmt.Println("entered function calcGaussandLapl printing shapes of img, " + fmt.Sprint(img.Size()) + ", gauss: " + fmt.Sprint(gauss.Size()) + ", lapl: " + fmt.Sprint(lapl.Size()))

	//   (img.Rows())/2, (img.Cols())/2

	gocv.PyrDown(*img, gauss, MakePyrKernel(*img, false), gocv.BorderDefault)
	gocv.PyrUp(*gauss, &temp, image.Pt(gauss.Rows()*2, gauss.Cols()*2), gocv.BorderDefault)

	//fmt.Println("went Pyr up and down shapes of img, " + fmt.Sprint(img.Size()) + ", gauss: " + fmt.Sprint(gauss.Size()) + ", lapl: " + fmt.Sprint(lapl.Size()) + ", temp: " + fmt.Sprint(temp.Size()))
	//fmt.Println("try to subtract temp, " + fmt.Sprint(temp.Size()) + ", from img: " + fmt.Sprint(img.Size()))
	gocv.Subtract(*img, temp, lapl)
	return nil
}

func ImageTo64float(img gocv.Mat) gocv.Mat {
	img.ConvertTo(&img, gocv.MatTypeCV64FC3)
	return img
}

func ImageTo8Int(img gocv.Mat) gocv.Mat {
	img.ConvertTo(&img, gocv.MatTypeCV8UC3)
	return img
}

func CreatePyramid(Pyr *Pyramid, img gocv.Mat, levels int) {

	*Pyr = make([]gocv.Mat, (levels + 1))

	gaussian_img := gocv.NewMat()
	defer gaussian_img.Close()

	laplacian_img := gocv.NewMat()
	defer laplacian_img.Close()

	//fmt.Println(img.Size())
	//fmt.Println(len(*Pyr))

	for j := range *Pyr {

		//fmt.Println("Vorgang nr: " + fmt.Sprint(j))

		CalcGaussandLapl(&img, &gaussian_img, &laplacian_img)

		if j == len(*Pyr)-1 {
			//fmt.Println("entered if")
			laplacian_img = img.Clone()
		} else {
			//fmt.Println("entered else")
			img = gaussian_img.Clone()
		}

		(*Pyr)[j] = laplacian_img.Clone()
	}
}

func ReconstructImageFromPyramid(Pyr Pyramid) *gocv.Mat {
	Result := Pyr[len(Pyr)-1].Clone()

	for i := len(Pyr) - 2; i >= 0; i-- {
		gocv.PyrUp(Result, &Result, image.Pt(Result.Rows()*2, Result.Cols()*2), gocv.BorderDefault)
		gocv.Add(Result, Pyr[i], &Result)
		//fmt.Println("Added Pyr with shape (" + fmt.Sprint(Result.Size()) + ") and Result with shape (" + fmt.Sprint(Result.Size()) + ")")
	}
	return &Result
}

func CreateImagesFromTimelines(TP [][]TimeLine) []gocv.Mat {

	return nil
}
