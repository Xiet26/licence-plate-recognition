package main

import (
	"flag"
	"fmt"
	"licence-plate-recognition/utilities"
	"log"
	"math/rand"
	"os"
	"runtime"
)

var (
	trainFile string
	testFile  string

	percentTrain float64
	percentTest  float64

	source string
)

func main() {
	flag.Parse()

	if percentTrain > 100 {
		log.Fatal("percent is not larger than 100")
	}
	percentTest = 100 - percentTrain

	allFolders, e := utilities.ListFolders(source)
	if e != nil || len(allFolders) == 0 {
		log.Fatal("no data")
	}

	allFiles := make([]string, 0)
	for _, v := range allFolders {
		files, e := utilities.ListFiles(v)
		if e != nil {
			continue
		}
		allFiles = append(allFiles, files...)
	}

	train, e := os.Create(trainFile)
	if e != nil {
		log.Fatal(e)
	}
	defer train.Close()

	test, e := os.Create(testFile)
	if e != nil {
		log.Fatal(e)
	}
	defer test.Close()

	existed := make(map[int]bool)

	countTest := 0
	for {
		if float64(countTest)/float64(len(allFiles))*100 > percentTest {
			break
		}

		i := rand.Intn(len(allFiles))
		if existed[i] {
			continue
		}

		existed[i] = true

		line, err := utilities.ToLineCSV(allFiles[i])
		if err != nil {
			continue
		}

		_, e = test.WriteString(fmt.Sprintf("%s\n", line))
		if e != nil {
			continue
		}

		countTest++
		if countTest%100 == 0 {
			fmt.Println("test: ", countTest)
		}
	}

	for k, v := range allFiles {
		if existed[k] {
			continue
		}

		line, err := utilities.ToLineCSV(v)
		if err != nil {
			continue
		}

		_, e = train.WriteString(fmt.Sprintf("%s\n", line))
		if e != nil {
			continue
		}

		if k%100 == 0 {
			fmt.Println("train: ", k)
		}
	}

	fmt.Println("done")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&trainFile, "trainFile", "train.csv", "data for train")
	flag.StringVar(&testFile, "testFile", "test.csv", "data for test")

	flag.Float64Var(&percentTrain, "percentTrain", 70, "percent data for train - default 70")

	flag.StringVar(&source, "source", "dataset", "folder contain data")
}
