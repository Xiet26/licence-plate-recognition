package main

import (
	"fmt"
	"licence-plate-recognition/licencePlateReconition"
)

var (
	lengthOfLicencePlate = 8
	filePath             = "../licenceplatesimage/51A99999.jpg" //image name
	modelPath            = "train.svm"
)

func main() {
	result, err := licencePlateReconition.Regconize(filePath, modelPath, true, lengthOfLicencePlate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
