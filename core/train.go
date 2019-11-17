package main

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
	"runtime"
)

func main() {
	param := libSvm.NewParameter()
	param.SvmType = libSvm.C_SVC
	param.KernelType = libSvm.LINEAR
	param.C = 0.5

	problem, e := libSvm.NewProblem("/home/phuoc/work-go/src/git.cyradar.com/phuocnn/licence-plate-recognition/train.csv", param)
	if e != nil {
		panic(e)
	}

	model := libSvm.NewModel(param)

	e = model.Train(problem)
	if e != nil {
		panic(e)
	}

	e = model.Dump("train.svm")
	if e != nil {
		panic(e)
	}

	fmt.Println("Trained!")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}