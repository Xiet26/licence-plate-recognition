package main

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
	"image"
	"licence-plate-recognition/utilities"
	"sort"
)

const (
	FOLDER_PART = `licenceplatesimage` //image folder
	EXPAND_SIZE = 4
)

var (
	lengthOfLicencePlate = 8
	fileName             = "51A99999.jpg" //image name
	dataMap              = utilities.DataMap()
)

func main() {
	//read image
	filePath := fmt.Sprintf("%s/%s", FOLDER_PART, fileName)
	fmt.Println(filePath)
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("empty image")
		return
	}
	//________________Car image__________________
	//Pre-Processing
	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	bilateral := gocv.NewMat()
	gocv.BilateralFilter(grayImg, &bilateral, 9, 75, 75)

	equal := gocv.NewMat()
	gocv.EqualizeHist(bilateral, &equal)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{5, 5})
	defer kernel.Close()
	morph := gocv.NewMat()
	gocv.MorphologyExWithParams(equal, &morph, gocv.MorphOpen, kernel, 20, gocv.BorderDefault)

	sub := gocv.NewMat()
	gocv.Subtract(equal, morph, &sub)

	thresh := gocv.NewMat()
	gocv.Threshold(sub, &thresh, 0, 255, gocv.ThresholdOtsu)

	canny := gocv.NewMat()
	gocv.Canny(thresh, &canny, 250, 255)

	kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Point{3, 3})
	dilate := gocv.NewMat()
	gocv.Dilate(canny, &dilate, kernel)

	//Find licence plate
	//Find contour
	contours := gocv.FindContours(dilate, gocv.RetrievalTree, gocv.ChainApproxSimple)

	//Sort contour(Area)
	sort.Slice(contours, func(i, j int) bool {
		return gocv.ContourArea(contours[i]) > gocv.ContourArea(contours[j])
	})

	//Find and border licence plate
	var licencePlateBound image.Rectangle
	for _, contour := range contours {
		peri := gocv.ArcLength(contour, true)
		approx := gocv.ApproxPolyDP(contour, 0.06*peri, true)
		if len(approx) == 4 {
			min := approx[0]
			max := approx[2]
			min.X -= EXPAND_SIZE
			min.Y -= EXPAND_SIZE
			max.X += EXPAND_SIZE
			max.Y += EXPAND_SIZE
			licencePlateBound = gocv.BoundingRect(approx)
			//licencePlateBound = image.Rectangle{Min:min, Max:max}
			break
		}
	}
	gocv.Rectangle(&img, licencePlateBound, colornames.Black, 6)
	//cut licence plate
	licencePlate := img.Region(licencePlateBound)

	//________________Licence plate image__________________
	//Pre-Processing
	licencePlateGray := gocv.NewMat()
	gocv.CvtColor(licencePlate, &licencePlateGray, gocv.ColorBGRToGray)

	licencePlateBlur := gocv.NewMat()
	gocv.GaussianBlur(licencePlateGray, &licencePlateBlur, image.Point{
		X: 3,
		Y: 3,
	}, 0, 0, gocv.BorderDefault)

	licencePlateThresh := gocv.NewMat()
	gocv.Threshold(licencePlateBlur, &licencePlateThresh, 120, 255, gocv.ThresholdBinaryInv)

	kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Point{
		X: 3,
		Y: 3,
	})
	licencePlateMor := gocv.NewMat()
	gocv.MorphologyEx(licencePlateThresh, &licencePlateMor, gocv.MorphDilate, kernel)

	//Find Contour
	contourNumbers := gocv.FindContours(licencePlateMor, gocv.RetrievalList, gocv.ChainApproxSimple)

	//Sort contour(area)
	contourNumberAreas := make(map[float64][]image.Point)
	var keyNumbers []float64
	for _, v := range contourNumbers {
		key := gocv.ContourArea(v)
		contourNumberAreas[key] = v
		keyNumbers = append(keyNumbers, key)
	}
	sort.Float64s(keyNumbers)

	// this here
	//Find number in licence plate
	contourNumbers = [][]image.Point{}
	for i := len(keyNumbers); i > 0; i-- {
		if i > len(keyNumbers)-lengthOfLicencePlate-2 && i < len(keyNumbers)-1 {
			contourNumbers = append(contourNumbers, contourNumberAreas[keyNumbers[i-1]])
		}
	}

	characterImg := make([]gocv.Mat, 0)
	//Border number in licence plate
	for _, v := range contourNumbers {
		rect := gocv.BoundingRect(v)
		//gocv.Rectangle(&licencePlate, rect, colornames.Red, 0)
		//expand border
		min := rect.Min
		max := rect.Max
		min.X -= EXPAND_SIZE
		min.Y -= EXPAND_SIZE
		max.X += EXPAND_SIZE
		max.Y += EXPAND_SIZE

		rectNew := image.Rectangle{min, max}
		characterImg = append(characterImg, licencePlate.Region(rectNew))
	}

	imgForRecognize := make(map[int]gocv.Mat)
	for k, v := range characterImg {
		tmp := gocv.NewMat()
		gocv.Resize(v, &tmp, image.Point{28, 28}, 0, 0, gocv.InterpolationNearestNeighbor)
		imgForRecognize[k] = tmp
	}
	ShowImg(licencePlate, "gray")

	//load model
	model := libSvm.NewModelFromFile("train.svm")

	for i := 0; i < len(imgForRecognize); i++ {
		grayImg := gocv.NewMat()
		gocv.CvtColor(imgForRecognize[i], &grayImg, gocv.ColorBGRToGray)
		ShowImg(grayImg, "gray")

		thresh := gocv.NewMat()
		gocv.Threshold(grayImg, &thresh, 128, 255, gocv.ThresholdBinaryInv)

		gocv.IMWrite(fmt.Sprintf("%d.png", i), thresh)

		x := utilities.MatToMapIntFloat(thresh)

		data := model.Predict(x)
		fmt.Println(fmt.Sprintf("%d: %s", i, dataMap[data]))

		grayImg.Close()
		thresh.Close()
	}

	//Show result
	//ShowImg(grayImg, "gray")
	//ShowImg(bilateral, "bilateral")
	//ShowImg(equal, "equal")
	//ShowImg(morph, "morph")
	//ShowImg(sub, "sub")
	//ShowImg(thresh, "thresh")
	//ShowImg(canny, "canny")
	//ShowImg(dilate, "dilate")
	//ShowImg(img, "Result detect")
}

func ShowImg(img gocv.Mat, name string) {
	window := gocv.NewWindow(name)
	defer window.Close()
	for {
		window.ResizeWindow(800, 600)
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
