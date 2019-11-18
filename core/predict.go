package main

import (
	"fmt"
	"git.cyradar.com/phuocnn/licence-plate-recognition/core/utilities"
	libSvm "github.com/ewalker544/libsvm-go"
)

func main() {
	model := libSvm.NewModelFromFile("/home/phuoc/work-go/src/git.cyradar.com/phuocnn/licence-plate-recognition/train.svm")

	filename := "/home/phuoc/work-go/src/git.cyradar.com/phuocnn/licence-plate-recognition/dataset/A/A_767.jpg"
	x := utilities.ToMapIntFloat(filename)

	result := model.Predict(x)

	dataMap := utilities.DataMap()
	fmt.Println("result: ", dataMap[result])
}
