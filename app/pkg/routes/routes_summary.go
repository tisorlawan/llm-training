package routes

import (
	"app/pkg"
	"app/pkg/api"
	"app/pkg/config"
	"app/templates"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RoutesSummary(r *gin.Engine, prefix string) {
	route := r.Group(prefix)
	summarizeSse := make(chan pkg.StatusMessage)

	route.GET("/summary", func(ctx *gin.Context) {
		templ.Handler(templates.Summarize()).ServeHTTP(ctx.Writer, ctx.Request)
	})

	route.GET("/summary.ws", func(ctx *gin.Context) {
		w, r := ctx.Writer, ctx.Request
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer c.Close()
		for token := range summarizeSse {
			c.WriteJSON(token)
		}
	})

	route.POST("/summary", func(ctx *gin.Context) {
		var chatReq pkg.SummaryRequest
		err := ctx.Bind(&chatReq)
		if err != nil {
			log.Println("[error]", err)
		}

		c := make(chan pkg.StatusMessage)
		summarizeSse = c

		model := "whisper.cpp/models/ggml-base.bin"

		go func() {
			hashBytes := md5.Sum([]byte(chatReq.URL + model))
			hash := hex.EncodeToString(hashBytes[:])
			filePath := fmt.Sprintf("content/%s.wav", hash)

			// download audio
			if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
				c <- pkg.StatusMessage{Type: "status", Data: "Downloading video ..."}
				cmd := exec.Command(
					"yt-dlp",
					"--extract-audio",
					"--audio-format", "wav",
					"--postprocessor-args", "-ar 16000",
					"-o", filePath,
					chatReq.URL,
				)
				cmd.Stdout = os.Stdout

				if err := cmd.Run(); err != nil {
					fmt.Println("could not run command: ", err)
				}
			}

			// transcribe audio
			lrcFile := fmt.Sprintf("content/%s.lrc", hash)
			if _, err := os.Stat(lrcFile); errors.Is(err, os.ErrNotExist) {
				c <- pkg.StatusMessage{Type: "status", Data: "Transcribing video ..."}
				cmd := exec.Command(
					"./whisper.cpp/main",
					"-m", model,
					"-f", filePath,
					"-of", fmt.Sprintf("content/%s", hash),
					"--language", "id",
					"--output-lrc",
				)
				cmd.Stdout = os.Stdout

				if err := cmd.Run(); err != nil {
					fmt.Println("could not run command: ", err)
				}
			}
			contentBytes, _ := os.ReadFile(lrcFile)
			content := string(contentBytes)

			cfgAny, _ := ctx.Get("config")
			cfg, _ := cfgAny.(*config.Configuration)

			tokenStream, _ := api.SummarizeTranscript(
				c,
				*cfg.OllamaModel,
				content,
				api.NewChatOptions().WithTemperature(0.5).WithNumPredict(1000),
			)

			for token := range tokenStream {
				c <- pkg.StatusMessage{Type: "token", Data: token}
			}
		}()
		templ.Handler(templates.SummarizeSSE()).ServeHTTP(ctx.Writer, ctx.Request)
	})
}
