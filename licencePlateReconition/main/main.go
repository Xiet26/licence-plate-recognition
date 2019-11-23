package main

import (
	"flag"
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"licence-plate-recognition/licencePlateReconition"
	"licence-plate-recognition/utilities"
	"log"
	"runtime"
)

var (
	lengthOfLicencePlate int
	filename             string
	modelPath            string
)

func main() {
	flag.Parse()

	fmt.Println("Loading model...")
	model := libSvm.NewModelFromFile(modelPath)
	dataMap := utilities.DataMap()

	fmt.Println("Licence Plate Detect...")
	licencePlateDetected, e := licencePlateReconition.Detect(filename, true, lengthOfLicencePlate)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("Licence Plate Recognize...")
	result, e := licencePlateReconition.Recognize(licencePlateDetected, model, dataMap)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("Result: ", result)

}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&filename, "filename", "../licenceplatesimage/51A99999.jpg", "image for recognize")
	flag.StringVar(&modelPath, "model", "train.svm", "model for svm")
	flag.IntVar(&lengthOfLicencePlate, "length", 8, "length of licence plate")
}
