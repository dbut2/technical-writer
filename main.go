// Package main is the entry point for the Technical Writer application.
// This program processes project files and leverages the OpenAI API to
// automatically generate documentation for those files.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// instruction provides guidance for the Technical Writer application on how
// to interact with files and respond to the OpenAI API.
const instruction = `You are a technical writer. 

Your job is to supply package level documentation that would ease an external developer to understand the program. You might also suggest creating a README.md file if you think it is necessary.

You will receive 1 message per file, where the message contains the file contains inside a code block titled with the filename.

If you think a change should be made to the file you should reply with the entire file with the documentation added, and if you want to create a file you should add a new code block with that file titled with the new files filename. If you do not make any changes to a file you may omit it from your response. Do not omit code when editing a block. Do not add any other message outside of the code blocks. Reply with 1 file per message.

If you do not wish to edit any more files, simply reply "STOP".`

// main is the entry point of the application.
func main() {
	openaiToken := os.Getenv("OPENAI_API_KEY") // Load the OpenAI API key from environment variables.

	client := openai.NewClient(openaiToken) // Create a new OpenAI client.
	_ = client

	files, err := listAllFiles("/github/workspace") // List all files in the specified directory.
	if err != nil {
		panic(err.Error()) // Handle potential errors during file listing.
	}

	err = document(context.Background(), client, files) // Begin the documentation process for each file.
	if err != nil {
		panic(err.Error()) // Handle potential errors during the documentation process.
	}
}

// listAllFiles recursively lists all files within the given directory, while
// avoiding version control directories that start with ".".
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

	allowList := strings.Split(os.Getenv("ALLOW_LIST"), ",")
	denyList := strings.Split(os.Getenv("DENY_LIST"), ",")

	var allowedFiles []string
	for _, file := range files {
		allowed := false

		for _, allow := range allowList {
			if strings.Contains(file, allow) {
				allowed = true
				break
			}
		}

		for _, deny := range denyList {
			if strings.Contains(file, deny) {
				allowed = false
				break
			}
		}

		if !allowed {
			continue
		}

		allowedFiles = append(allowedFiles, file)
	}

	return allowedFiles, nil
}

// document interacts with the OpenAI API to generate and apply
// documentation to each file in the provided list.
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

		req.Messages = append(req.Messages, openai.ChatCompletionMessage{ // Append a message for each file.
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("```%s\n%s\n```", file, contents),
		})
	}

	for {
		chat, err := client.CreateChatCompletion(ctx, req) // Send a chat completion request to OpenAI.
		if err != nil {
			return err
		}
		req.Messages = append(req.Messages, chat.Choices[0].Message)

		if chat.Choices[0].Message.Content == "STOP" {
			break
		}

		file, contents := parseResponse(chat.Choices[0].Message.Content)

		err = os.WriteFile(file, contents, 0644) // Write the updated documentation back to the file.
		if err != nil {
			return err
		}
	}

	return nil
}

// parseResponse parses a documentation response from the OpenAI API
// and extracts the filename and content to be applied to the file.
func parseResponse(response string) (string, []byte) {
	lines := strings.Split(response, "\n")
	file := strings.Trim(lines[0], "`")
	contents := []byte(strings.Join(lines[1:len(lines)-1], "\n"))

	return file, contents
}
