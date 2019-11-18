package main

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"licence-plate-recognition/utilities"
)

func main() {
	model := libSvm.NewModelFromFile("../train.svm")

	filename := "../dataset/A/A_767.jpg"
	x := utilities.ToMapIntFloat(filename)

	result := model.Predict(x)

	dataMap := utilities.DataMap()
	fmt.Println("result: ", dataMap[result])
}
