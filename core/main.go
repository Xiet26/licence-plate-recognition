package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
	"image"
	"sort"
)

const FOLDER_PART = `` //image folder

var (
	lengthOfLicencePlate = 9
	fileName             = `` //image name
)

func main() {
	//read image
	filePath := fmt.Sprintf("%s%s", FOLDER_PART, fileName)
	fmt.Println(filePath)
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("empty")
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
	contoursArea := make(map[float64][]image.Point)
	var keys []float64
	for _, v := range contours {
		key := gocv.ContourArea(v)
		contoursArea[key] = v
		keys = append(keys, key)
	}
	sort.Float64s(keys)

	contours = [][]image.Point{}
	for i := len(keys); i > 0; i-- {
		contours = append(contours, contoursArea[keys[i-1 ]])
	}

	//Find and border licence plate
	var licencePlateBound image.Rectangle
	for _, contour := range contours {
		peri := gocv.ArcLength(contour, true)
		approx := gocv.ApproxPolyDP(contour, 0.06*peri, true)
		if len(approx) == 4 {
			licencePlateBound = gocv.BoundingRect(approx)
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
	licencePlateThre := gocv.NewMat()
	gocv.Threshold(licencePlateBlur, &licencePlateThre, 120, 255, gocv.ThresholdBinaryInv)
	kernel = gocv.GetStructuringElement(gocv.MorphRect, image.Point{
		X: 3,
		Y: 3,
	})
	licencePlateMor := gocv.NewMat()
	gocv.MorphologyEx(licencePlateThre, &licencePlateMor, gocv.MorphDilate, kernel)
	//Find Contour
	contourNumbers := gocv.FindContours(licencePlateMor, gocv.RetrievalTree, gocv.ChainApproxSimple)
	//Sort contour(area)
	contourNumberAreas := make(map[float64][]image.Point)
	var keyNumbers []float64
	for _, v := range contourNumbers {
		key := gocv.ContourArea(v)
		contourNumberAreas[key] = v
		keyNumbers = append(keyNumbers, key)
	}
	sort.Float64s(keyNumbers)
	contourNumbers = [][]image.Point{}
	//Find number in licence plate
	for i := len(keyNumbers); i > 0; i-- {
		if i > len(keyNumbers)-lengthOfLicencePlate-1 && i < len(keyNumbers)-1 {
			contourNumbers = append(contourNumbers, contourNumberAreas[keyNumbers[i-1]])
		}
	}
	//Border number in licence plate
	for _, v := range contourNumbers {
		rect := gocv.BoundingRect(v)
		gocv.Rectangle(&licencePlate, rect, colornames.Red, 2)
	}
	//Show result
	window := gocv.NewWindow("hello")
	for {
		window.IMShow(licencePlate)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
