package main

import (
	"fmt"
	"os"
)

func main() {
	listAllFiles()
}

func listAllFiles() {
	file, err := os.Open("/github/workspace")
	if err != nil {
		panic(err.Error())
	}

	files, err := file.Readdir(0)
	if err != nil {
		panic(err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
