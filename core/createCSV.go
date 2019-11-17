package main

import (
	"fmt"
	"git.cyradar.com/phuocnn/licence-plate-recognition/core/utilities"
	"os"
)

func main() {
	src := "/home/phuoc/work-go/src/git.cyradar.com/phuocnn/licence-plate-recognition/dataset"

	allFolder := utilities.ListFolders(src)
	if len(allFolder) == 0 {
		fmt.Println("error no data")
		return
	}

	allFiles := make([]string, 0)
	for _, v := range allFolder {
		allFiles = append(allFiles, utilities.ListFiles(v)...)
	}

	// write data to csv
	f, e := os.Create("train.csv")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	for k, v := range allFiles {
		line := utilities.ToLineCSV(v)
		_, e := f.WriteString(fmt.Sprintf("%s\n", line))
		if e != nil {
			continue
		}

		if k%100 == 0 {
			fmt.Println(k)
		}
	}

}
