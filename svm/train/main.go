package main

import (
	"flag"
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"runtime"
)

var (
	trainFile  string
	resultFile string
)

func main() {
	param := libSvm.NewParameter()
	param.SvmType = libSvm.C_SVC
	param.KernelType = libSvm.LINEAR
	param.C = 0.5

	problem, e := libSvm.NewProblem(trainFile, param)
	if e != nil {
		panic(e)
	}

	model := libSvm.NewModel(param)

	e = model.Train(problem)
	if e != nil {
		panic(e)
	}

	e = model.Dump(resultFile)
	if e != nil {
		panic(e)
	}

	fmt.Println("Trained!")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&trainFile, "trainFile", "train.csv", "data for train")
	flag.StringVar(&resultFile, "resultFile", "train.svm", "file model after train")
}
