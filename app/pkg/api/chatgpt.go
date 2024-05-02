package api

import (
	"app/pkg"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

func init() {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		log.Fatalln("environment variable OPENAI_API_KEY must be provided")
	}
	openaiClient = openai.NewClient(openAIAPIKey)
}

func ChatGPTChatStream(messages []pkg.Message) (chan string, error) {
	c := make(chan string)
	go func() {
		stream, err := openaiClient.CreateChatCompletionStream(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: messagesToChatCompletionMessages(messages),
				Stream:   true,
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
		}
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(c)
				break
			}
			if err != nil {
				break
			}
			if len(resp.Choices) > 0 {
				diff := resp.Choices[0].Delta.Content
				c <- diff
			}

		}
	}()
	return c, nil
}

func ChatGPTCreateImageDalle3(prompt string) string {
	c, err := openaiClient.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Model:          "dall-e-3",
		N:              1,
		Size:           "1024x1024",
		Quality:        "standard",
		ResponseFormat: "b64_json",
	})
	if err != nil {
		log.Println("[error] request image", err)
	}
	return decodeAndSaveB64Image(c.Data[0].B64JSON)
}

func ChatGPTCreateImageDalle2(prompt string, size string) string {
	c, err := openaiClient.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Model:          "dall-e-2",
		N:              1,
		Size:           size,
		Quality:        "standard",
		ResponseFormat: "b64_json",
	})
	if err != nil {
		log.Println("[error] request image", err)
	}
	return decodeAndSaveB64Image(c.Data[0].B64JSON)
}

func ChatGPTCreateImageDalle2Stream(prompt string, size string) chan string {
	c := make(chan string)
	go func() {
		filename := ChatGPTCreateImageDalle2(prompt, size)
		c <- filename
	}()
	return c
}

func messagesToChatCompletionMessages(messages []pkg.Message) []openai.ChatCompletionMessage {
	cmpMessages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: pkg.DefaultSystemPrompt,
		},
	}

	for _, msg := range messages {
		role := "user"
		if !msg.IsUser {
			role = "assistant"
		}
		cmpMessages = append(cmpMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Text,
		})
	}
	return cmpMessages
}

func decodeAndSaveB64Image(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Println("[error] decoding image")
	}
	filename := fmt.Sprintf("images/image-%x.png", md5.Sum(decoded))
	os.WriteFile(filename, decoded, 0644)
	return filename
}
