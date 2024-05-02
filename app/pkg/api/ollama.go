package api

import (
	"app/pkg"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var OllamaURL string = os.Getenv("OLLAMA_URL")

func init() {
	if OllamaURL == "" {
		log.Fatalln("environment variable OLLAMA_URL must be provided")
	}
}

type Model struct {
	ID   string `json:"model"`
	Name string `json:"name"`
}

type Models struct {
	Models []Model `json:"models"`
}

type ChatOptions struct {
	NumPredict  *int     `json:"num_predict,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
}

func NewChatOptions() *ChatOptions {
	return &ChatOptions{}
}

func (c *ChatOptions) WithNumPredict(n int) *ChatOptions {
	c.NumPredict = &n
	return c
}

func (c *ChatOptions) WithTemperature(t float32) *ChatOptions {
	c.Temperature = &t
	return c
}

type ChatRequest struct {
	Model    string           `json:"model"`
	Messages []pkg.RawMessage `json:"messages"`
	Stream   bool             `json:"stream"`
	Options  *ChatOptions     `json:"options,omitempty"`
}

type ChatTokenStreamResponse struct {
	Done    bool           `json:"done"`
	Message pkg.RawMessage `json:"message"`
}

func SummarizeTranscript(statusChan chan pkg.StatusMessage, model string, transcript string, options *ChatOptions) (chan string, error) {
	statusChan <- pkg.StatusMessage{Type: "status", Data: "Summarizing transcript ..."}
	messages := pkg.MessagesToRawMessages([]pkg.Message{
		{
			Text:      transcript,
			IsUser:    true,
			Timestamp: time.Now(),
			AudioSrc:  "",
		},
	}, pkg.SummarizeTranscriptPromptExample)

	fmt.Printf("%v\n", messages)

	req := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
		Options:  options,
	}
	return streamResponse(req)
}

func ChatStreamOllama(
	model string,
	messages []pkg.Message,
	options *ChatOptions,
) (chan string, error) {
	req := ChatRequest{
		Model:    model,
		Messages: pkg.MessagesToRawMessages(messages, pkg.DefaultSystemPrompt),
		Stream:   true,
		Options:  options,
	}
	return streamResponse(req)
}

func streamResponse(req ChatRequest) (chan string, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("can't marshal request")
	}

	c := make(chan string)
	go func() {
		res, err := http.Post(
			OllamaURL+"/api/chat",
			"application/json",
			bytes.NewReader(data),
		)
		if err != nil {
			log.Println("[error] Chat /api/chat:", err)
		}

		idx := 0
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			bytes := scanner.Bytes()

			var jres ChatTokenStreamResponse
			if err = json.Unmarshal(bytes, &jres); err != nil {
				log.Println("[error] getModels: error parsing json", err)
			}

			if jres.Done {
				close(c)
			} else {
				if idx == 0 {
					c <- strings.TrimPrefix(strings.TrimPrefix(jres.Message.Content, " "), "\n")
				} else {
					c <- jres.Message.Content
				}
			}
			idx += 1
		}
	}()

	return c, nil
}

func ListModels() ([]Model, error) {
	res, err := http.Get(OllamaURL + "/api/tags")
	if err != nil {
		log.Println("[error] ListModels /api/tags:", err)
		return nil, err
	}
	body := res.Body
	dec := json.NewDecoder(body)

	var models Models
	if err = dec.Decode(&models); err != nil {
		log.Println("[error] getModels: error parsing json", err)
		return nil, err
	}
	return models.Models, nil
}
