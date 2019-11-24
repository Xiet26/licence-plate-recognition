package main

import (
	"flag"
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"log"
	"os"
	"runtime"
)

var (
	testFile   string
	modelFile  string
	reportFile string
)

func main() {
	flag.Parse()

	model := libSvm.NewModelFromFile(modelFile)

	tmpParam := libSvm.NewParameter()
	problems, e := libSvm.NewProblem(testFile, tmpParam)
	if e != nil {
		log.Fatal("wrong format data to test: ", e)
	}

	positive := 0
	negative := 0
	fileNegative := make(map[string]string)

	count := 0
	problems.Begin()
	for {
		if count%100 == 0 {
			fmt.Println(count)
		}

		if problems.Done() {
			break
		}

		expect, x := problems.GetLine()
		problems.Next()
		count++

		out := model.Predict(x)
		if expect == out {
			positive++
			continue
		}
		negative++
		fileNegative[string(uint64(expect))] = string(uint64(out))
	}

	report, e := os.Create(reportFile)
	if e != nil {
		panic(e)
	}
	defer report.Close()

	for k, v := range fileNegative {
		_, _ = report.WriteString(fmt.Sprintf("expect: %s - but get: %s\n", k, v))
	}

	_, _ = report.WriteString(fmt.Sprintf("result test: %f", float64(negative)/float64(negative+positive)))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&testFile, "testFile", "test.csv", "data for test")
	flag.StringVar(&modelFile, "modelFile", "train.svm", "model for test")
	flag.StringVar(&reportFile, "reportFile", "report.csv", "result after test")
}
