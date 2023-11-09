package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

const instruction = `You are a technical writer. You should supply code suggestions that increase readability for developers integrating with the code by creating comments, editing existing comments for readability and supply other suggestions that would help with developer experience.

You must reply with just the existing code edited. Don't add any other messages. Do not omit code. Do not reply in a code block.`

func main() {
	openaiToken := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(openaiToken)
	_ = client

	files, err := listAllFiles("/github/workspace")
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	eg := errgroup.Group{}
	for _, file := range files {
		file := file
		eg.Go(func() error {
			return document(ctx, client, file)
		})
	}
	if err = eg.Wait(); err != nil {
		panic(err.Error())
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
		fullPath := dir + "/" + fileInfo.Name()

		if fileInfo.IsDir() {
			if strings.Contains(fileInfo.Name(), ".") {
				continue
			}

			subFiles, err := listAllFiles(fullPath)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, fullPath)
		}
	}

	return files, nil
}

func document(ctx context.Context, client *openai.Client, file string) error {
	fmt.Println("Documenting " + file + "...")
	failed := "Documenting " + file + " failed"
	defer func() {
		fmt.Println(failed)
	}()

	contents, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	chat, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4-1106-preview",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: instruction,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(contents),
			},
		},
	})
	if err != nil {
		return err
	}

	newContents := chat.Choices[0].Message.Content

	err = os.WriteFile(file, []byte(newContents), 0644)
	if err != nil {
		return err
	}

	failed = "Documented " + file + " successfully"
	return nil
}
