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

func ToMapIntFloat(filename string) (map[int]float64, error) {
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

func DataMap() map[float64]string {
	result := make(map[float64]string)

	// 0 - 9
	for i := 48; i <= 57; i++ {
		result[float64(i)] = string(i)
	}

	// A - Z
	for i := 65; i <= 90; i++ {
		result[float64(i)] = string(i)
	}

	return result
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

func ListFolders(folder string) ([]string, error) {
	result := make([]string, 0)

	folders, e := ioutil.ReadDir(folder)
	if e != nil {
		return result, e
	}

	for _, v := range folders {
		if !v.IsDir() {
			continue
		}

		result = append(result, fmt.Sprintf("%s/%s", folder, v.Name()))
	}

	return result, nil
}

func ListFiles(folder string) ([]string, error) {
	result := make([]string, 0)

	files, e := ioutil.ReadDir(folder)
	if e != nil {
		return result, e
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		result = append(result, fmt.Sprintf("%s/%s", folder, v.Name()))
	}

	return result, nil
}

func CreateCSVFileFromData(src string) error{
	allFolder, err := ListFolders(src)
	if err != nil {
		return err
	}
	if len(allFolder) == 0 {
		return fmt.Errorf(ERROR_EMPTY_DATA)
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
	return nil
}

func ShowImg(img gocv.Mat, name string, width int, height int) {
	window := gocv.NewWindow(name)
	defer window.Close()
	for {
		window.ResizeWindow(width, height)
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}