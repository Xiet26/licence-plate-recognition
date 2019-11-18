package main

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"github.com/sirupsen/logrus"
	"licence-plate-recognition/utilities"
	"os"
)

var log *logrus.Logger
func main() {
	model := libSvm.NewModelFromFile("train.svm")
	filename := "../dataset/A/A_767.jpg"
	x, err := utilities.ToMapIntFloat(filename)
	if err != nil {
		log.Fatal(err)
	}

	result := model.Predict(x)

	dataMap, err := utilities.DataMap("../dataset")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result: ", dataMap[result])
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
	log = logrus.New()
	log.Out = os.Stdout
}