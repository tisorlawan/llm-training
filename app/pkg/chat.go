package pkg

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChatRequest struct {
	Message string `form:"message"`
}

type SummaryRequest struct {
	URL string `form:"url"`
}

type Message struct {
	Text      string
	IsUser    bool
	Timestamp time.Time
	AudioSrc  string
}

type RawMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StatusMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var DefaultSystemPrompt = `You are an AI Chatbot Assistant,
you should answer every question to the best of your knowledge. 
Answer [IMAGE] if you think the user asking for an image.
Always answer in Bahasa Indonesia.

For Example:

Question: When did Indonesia get their independence ?
Answer: 17 Agustus 1945

Question: Give me an image of Cat vs Dog!
Answer: [IMAGE]

Question: Generate an image of mars for me
Answer: [IMAGE]

Question: Buat sebuah gambar sepeda motor dikendarai seorang ninja
Answer: [IMAGE]`

var SummarizeTranscriptPromptExample = `
	You are an AI Chatbot Assistant. 
	Summarize the following audio transcript with lrc format. Answer in bullet points, wrap each point in HTML <p> with bullet point. 
	Answer in Bahasa Indonesia.
	If it is a news, summarize the news content.
	Don't add additional notes.
`

func MessagesToRawMessages(messages []Message, systemPrompt string) []RawMessage {
	rawMessages := []RawMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	for _, msg := range messages {
		role := "user"
		if !msg.IsUser {
			role = "assistant"
		}
		rawMessages = append(rawMessages, RawMessage{
			Role:    role,
			Content: msg.Text,
		})
	}

	fmt.Println("---------------------------------")
	rawJSON, _ := json.MarshalIndent(rawMessages, "", "  ")
	fmt.Println(string(rawJSON))

	return rawMessages
}
