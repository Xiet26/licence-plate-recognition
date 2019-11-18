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

func ToMapIntFloat(filename string) map[int]float64 {
	result := make(map[int]float64)
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("empty")
		return result
	}

	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	b := grayImg.ToBytes()

	for k, v := range b {
		result[k] = float64(v)
	}

	return result
}

func DataMap() map[float64]string {
	result := make(map[float64]string)
	folders, e := ioutil.ReadDir("/home/phuoc/work-go/src/git.cyradar.com/phuocnn/licence-plate-recognition/dataset")
	if e != nil {
		return result
	}

	for _, v := range folders {
		if !v.IsDir() {
			continue
		}

		result[float64(v.Name()[0])] = string(v.Name()[0])
	}

	return result
}

func ToLineCSV(filename string) string {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("empty")
		return ""
	}

	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	b := grayImg.ToBytes()

	line := fmt.Sprintf("%d ", filepath.Base(filename)[0])
	for k, v := range b {
		line = line + fmt.Sprintf("%d:%d ", k, v)
	}

	return line
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

func ListFiles(folder string) []string {
	result := make([]string, 0)

	files, e := ioutil.ReadDir(folder)
	if e != nil {
		return result
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

	return result
}

func CreateCSVFileFromData(src string) {
	allFolder := ListFolders(src)
	if len(allFolder) == 0 {
		fmt.Println("error no data")
		return
	}

	allFiles := make([]string, 0)
	for _, v := range allFolder {
		allFiles = append(allFiles, ListFiles(v)...)
	}

	// write data to csv
	f, e := os.Create("train.csv")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	for k, v := range allFiles {
		line := ToLineCSV(v)
		_, e := f.WriteString(fmt.Sprintf("%s\n", line))
		if e != nil {
			continue
		}
		if k%100 == 0 {
			fmt.Println(k)
		}
	}

}
