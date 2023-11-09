package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const instruction = `You are a technical writer. 

Your job is to supply package level comments that would ease an external developer to understand the program. You might also suggest creating a README.md file if you think it is necessary.

You will receive 1 message per file, where the message contains the file contains inside a code block titled with the filename.

If you think a change should be made to the file you should reply with the entire file returned with comments added, and if you want to create a file you should add a new code block with that file titled with the new files filename. Do not omit code when editing a block. Do not add any other message outside of the code blocks.`

func main() {
	openaiToken := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(openaiToken)
	_ = client

	files, err := listAllFiles("/github/workspace")
	if err != nil {
		panic(err.Error())
	}

	err = document(context.Background(), client, files)
	if err != nil {
		panic(err.Error())
	}
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
			if strings.HasPrefix(fileInfo.Name(), ".") {
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

func document(ctx context.Context, client *openai.Client, files []string) error {

	req := openai.ChatCompletionRequest{
		Model: "gpt-4-1106-preview",
		Messages: append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: instruction,
			},
		}),
	}

	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("```%s\n%s\n```", file, contents),
		})
	}

	for i := 0; i < 10; i++ {
		chat, err := client.CreateChatCompletion(ctx, req)
		if err != nil {
			return err
		}

		fmt.Println(chat.Choices[0])

		//err = os.WriteFile(file, []byte(newContents), 0644)
		//if err != nil {
		//	return err
		//}
	}

	return nil
}
