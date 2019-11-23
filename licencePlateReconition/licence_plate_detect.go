package licencePlateReconition

import (
	"fmt"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
	"image"
	"licence-plate-recognition/utilities"
	"sort"
)

const (
	EXPAND_SIZE = 4
)

func Detect(imagePath string, showImage bool, lengthOfLicencePlate int) (map[int]gocv.Mat, error) {
	//read image
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("empty image")
		return nil, fmt.Errorf(utilities.ERROR_EMPTY_IMAGE)
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

	grayImg.Close()
	bilateral.Close()
	equal.Close()
	morph.Close()
	sub.Close()
	morph.Close()
	sub.Close()
	thresh.Close()
	canny.Close()
	dilate.Close()
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

	//Find number in licence plate
	contourNumbers = [][]image.Point{}
	for i := len(keyNumbers); i > 0; i-- {
		if i > len(keyNumbers)-lengthOfLicencePlate-2 && i < len(keyNumbers)-1 {
			contourNumbers = append(contourNumbers, contourNumberAreas[keyNumbers[i-1]])
		}
	}

	var rects []image.Rectangle
	characterImg := make([]gocv.Mat, 0)
	//Border number in licence plate
	for _, v := range contourNumbers {
		rect := gocv.BoundingRect(v)
		//expand border
		min := rect.Min
		max := rect.Max
		min.X -= EXPAND_SIZE
		min.Y -= EXPAND_SIZE
		max.X += EXPAND_SIZE
		max.Y += EXPAND_SIZE

		rectNew := image.Rectangle{Min: min, Max: max}
		rects = append(rects, rectNew)
	}

	sort.Slice(rects, func(i, j int) bool {
		return rects[i].Min.X < rects[j].Min.X
	})
	for _, v := range rects {
		characterImg = append(characterImg, licencePlate.Region(v))
	}

	imgForRecognize := make(map[int]gocv.Mat)
	for k, v := range characterImg {
		tmp := gocv.NewMat()
		gocv.Resize(v, &tmp, image.Point{28, 28}, 0, 0, gocv.InterpolationNearestNeighbor)
		imgForRecognize[k] = tmp
	}
	if showImage {
		utilities.ShowImg(licencePlate, "Licence Plate", 1000, 500)
	}

	return imgForRecognize, nil
}
