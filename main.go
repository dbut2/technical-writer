package main

import (
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	args := os.Args
	if len(args) != 1 {
		panic("incorrect args passed, should be 1")
	}

	openaiToken := args[0]

	client := openai.NewClient(openaiToken)
	_ = client

	files, err := listAllFiles("/github/workspace")
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		err = document(file)
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println(files)
}

func listAllFiles(dir string) ([]string, error) {
	var files []string

	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfos, err := file.ReadDir(0)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subFiles, err := listAllFiles(fileInfo.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, dir+"/"+fileInfo.Name())
		}
	}

	return files, nil
}

func document(file string) error {
	fmt.Println(file)
	return nil
}
