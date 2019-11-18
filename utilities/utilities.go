package utilities

import (
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MatToMapIntFloat(mat gocv.Mat) map[int]float64 {
	result := make(map[int]float64)

	b := mat.ToBytes()

	for k, v := range b {
		result[k] = float64(v)
	}

	return result
}

func ToMapIntFloat(filename string) (map[int]float64, error ) {
	result := make(map[int]float64)
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		return result, fmt.Errorf(ERROR_EMPTY_IMAGE)
	}

	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	b := grayImg.ToBytes()

	for k, v := range b {
		result[k] = float64(v)
	}

	return result, nil
}

func DataMap(path string) (map[float64]string, error) {
	result := make(map[float64]string)
	folders, e := ioutil.ReadDir(path)
	if e != nil {
		return result, e
	}

	for _, v := range folders {
		if !v.IsDir() {
			continue
		}

		result[float64(v.Name()[0])] = string(v.Name()[0])
	}

	return result, nil
}

func ToLineCSV(filename string) (string, error) {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		return "", fmt.Errorf(ERROR_EMPTY_IMAGE)
	}

	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	b := grayImg.ToBytes()

	line := fmt.Sprintf("%d ", filepath.Base(filename)[0])
	for k, v := range b {
		line = line + fmt.Sprintf("%d:%d ", k, v)
	}

	return line, nil
}

func ListFolders(folder string) []string {
	result := make([]string, 0)

	folders, e := ioutil.ReadDir(folder)
	if e != nil {
		return result
	}

	for _, v := range folders {
		if !v.IsDir() {
			continue
		}

		result = append(result, fmt.Sprintf("%s/%s", folder, v.Name()))
	}

	return result
}

func ListFiles(folder string) ([]string, error) {
	result := make([]string, 0)

	files, e := ioutil.ReadDir(folder)
	if e != nil {
		return result, e
	}

	count := 0
	for _, v := range files {
		if v.IsDir() {
			continue
		}

		result = append(result, fmt.Sprintf("%s/%s", folder, v.Name()))
		count++
		if count > 500 {
			break
		}
	}

	return result, nil
}

func CreateCSVFileFromData(src string) {
	allFolder := ListFolders(src)
	if len(allFolder) == 0 {
		fmt.Println("error no data")
		return
	}

	allFiles := make([]string, 0)
	for _, v := range allFolder {
		listFiles, err := ListFiles(v)
		if err != nil {
			continue
		}
		allFiles = append(allFiles, listFiles...)
	}

	// write data to csv
	f, e := os.Create("train.csv")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	for k, v := range allFiles {
		line, err := ToLineCSV(v)
		if err != nil {
			continue
		}
		_, err = f.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			continue
		}
		if k%100 == 0 {
			fmt.Println(k)
		}
	}

}
