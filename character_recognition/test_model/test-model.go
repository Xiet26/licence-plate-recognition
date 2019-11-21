package main

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"licence-plate-recognition/utilities"
	"os"
	"path/filepath"
)

var log *logrus.Logger

func main() {
	model := libSvm.NewModelFromFile("train.svm")

	src := "../dataset"

	allFolder, err := utilities.ListFolders(src)
	if err != nil {
		log.Errorln(err)
	}
	if len(allFolder) == 0 {
		log.Errorln("error no data for test")
		return
	}

	allFiles := make([]string, 0)
	for _, v := range allFolder {
		files, err := ioutil.ReadDir(v)
		if err != nil {
			continue
		}

		count := 0
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			count++
			if count < 500 {
				continue
			}
			allFiles = append(allFiles, fmt.Sprintf("%s/%s", v, file.Name()))
		}

	}

	positive := 0
	negative := 0
	fileNegative := make(map[string]string)
	count := 0
	for _, v := range allFiles {
		if count % 100 == 0{
			fmt.Println(count)
		}
		count++
		x, err := utilities.ToMapIntFloat(v)
		if err != nil {
			log.Fatal(err)
		}

		expect := float64(filepath.Base(v)[0])

		out := model.Predict(x)
		if out == expect {
			positive++
			continue
		}

		negative++
		fileNegative[v] = string(uint64(out))
	}

	report, e := os.Create("test-model.csv")
	if e != nil {
		panic(e)
	}
	defer report.Close()

	for k, v := range fileNegative {
		_, _ = report.WriteString(fmt.Sprintf("%s - out: %s\n", k, v))
	}

	_, _ = report.WriteString(fmt.Sprintf("result test: %f", float64(negative)/float64(negative+positive)))
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
	log = logrus.New()
	log.Out = os.Stdout
}
