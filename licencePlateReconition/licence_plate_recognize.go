package licencePlateReconition

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"gocv.io/x/gocv"
	"licence-plate-recognition/utilities"
)

func Recognize(imgForRecognize map[int]gocv.Mat, model *libSvm.Model, dataMap map[float64]string) (string, error) {
	if model.NrClass() == 0 {
		return "", fmt.Errorf("no model")
	}

	var result string

	for i := 0; i < len(imgForRecognize); i++ {
		grayImg := gocv.NewMat()
		gocv.CvtColor(imgForRecognize[i], &grayImg, gocv.ColorBGRToGray)

		utilities.ShowImg(grayImg, fmt.Sprintf("%d", i+1), 50, 50)

		thresh := gocv.NewMat()
		gocv.Threshold(grayImg, &thresh, 128, 255, gocv.ThresholdBinaryInv)

		x := utilities.MatToMapIntFloat(thresh)

		data := model.Predict(x)

		result += dataMap[data]

		grayImg.Close()
		thresh.Close()
	}

	return result, nil
}
