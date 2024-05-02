package routes

import (
	"app/pkg"
	"app/pkg/api"
	"app/pkg/config"
	"app/templates"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func RoutesChat(r *gin.Engine, prefix string) {
	currentImageID := ""
	oldImageID := ""

	messages := []pkg.Message{}
	answers := map[string]chan string{}
	questions := map[string]string{}
	audios := map[string]string{}
	replies := map[string]string{}

	route := r.Group(prefix)

	route.GET("/", func(ctx *gin.Context) {
		cfgAny, _ := ctx.Get("config")
		cfg, _ := cfgAny.(*config.Configuration)
		if cfg.ChatModel == "ollama" {
			models, err := api.ListModels()
			if err != nil {
				log.Println("[error]", err)
			}
			templ.Handler(templates.Home(cfg, models, messages, currentImageID)).
				ServeHTTP(ctx.Writer, ctx.Request)
		} else {
			templ.Handler(templates.Home(cfg, []api.Model{}, messages, currentImageID)).ServeHTTP(ctx.Writer, ctx.Request)
		}
	})

	route.POST("/configuration", func(ctx *gin.Context) {
		time.Sleep(500 * time.Millisecond)
		var newConfig config.Configuration
		if err := ctx.Bind(&newConfig); err != nil {
			log.Println("[error]", err)
		}
		cfgAny, _ := ctx.Get("config")
		cfg, _ := cfgAny.(*config.Configuration)

		cfg.ChatModel = newConfig.ChatModel
		switch cfg.ChatModel {
		case "chatgpt":
			cfg.OpenaiAPIKey = newConfig.OpenaiAPIKey
		case "ollama":
			cfg.OllamaModel = newConfig.OllamaModel
		}
		cfg.ImageModel = newConfig.ImageModel
		cfg.UseAudio = newConfig.UseAudio
	})

	route.GET("/configuration", func(ctx *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		chatModel := ctx.Request.URL.Query().Get("chat_model")
		cfgAny, _ := ctx.Get("config")
		cfg, _ := cfgAny.(*config.Configuration)
		if chatModel == "ollama" {
			models, err := api.ListModels()
			if err != nil {
				log.Println("[error]", err)
			}
			templ.Handler(templates.OllamaConfiguration(models, cfg)).
				ServeHTTP(ctx.Writer, ctx.Request)
		} else {
			templ.Handler(templates.ChatGPTConfiguration(*cfg.OpenaiAPIKey)).ServeHTTP(ctx.Writer, ctx.Request)
		}
	})

	route.GET("/current-image", func(ctx *gin.Context) {
		if strings.HasPrefix(currentImageID, "images/") {
			templ.Handler(templates.ImageReady(currentImageID, oldImageID)).
				ServeHTTP(ctx.Writer, ctx.Request)
		} else {
			templ.Handler(templates.ImageLoading()).ServeHTTP(ctx.Writer, ctx.Request)
		}
	})

	route.GET("/audio/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		templ.Handler(templates.ChatWithAudio(audios[id], replies[id], id)).
			ServeHTTP(ctx.Writer, ctx.Request)
	})

	route.GET("/answer/:id", func(ctx *gin.Context) {
		cfgAny, _ := ctx.Get("config")
		cfg, _ := cfgAny.(*config.Configuration)

		w, r := ctx.Writer, ctx.Request
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer c.Close()

		reply := ""
		start := time.Now()
		id := ctx.Param("id")
		cx := answers[ctx.Param("id")]
		isReserved := false

		for s := range cx {
			if reply == "[IMAGE]" {
				break
			} else if len(reply) > len("[IMAGE]") {
				if !isReserved {
					if !cfg.UseAudio {
						c.WriteMessage(websocket.TextMessage, []byte(reply))
					}
					isReserved = true
				}
				if !cfg.UseAudio {
					c.WriteMessage(websocket.TextMessage, []byte(s))
				}
			}
			reply += s
		}

		if reply == "[IMAGE]" {
			imageID := "image-" + strconv.Itoa(int(time.Now().Unix()))
			currentImageID = imageID
			oldImageID = id
			c.WriteMessage(websocket.TextMessage, []byte("[IMAGE]_"+imageID))

			go func() {
				if cfg.ImageModel == "dall-e-3" {
					currentImageID = api.ChatGPTCreateImageDalle3(questions[id])
				} else {
					currentImageID = api.ChatGPTCreateImageDalle2(questions[id], "256x256")
				}
			}()

		} else if len(reply) < len("[IMAGE]") {
			if !cfg.UseAudio {
				c.WriteMessage(websocket.TextMessage, []byte(reply))
			}
		}

		replies[id] = reply

		audioSrc := ""
		if cfg.UseAudio && reply != "[IMAGE]" {
			fname, err := api.ProsaAudio(reply[:min(200, len(reply))], cfg.ProsaTTSModel)
			if err != nil {
				log.Println("[error]", err)
			}
			audios[id] = fname
			audioSrc = fname
			c.WriteMessage(websocket.TextMessage, []byte("[AUDIO]_"+id))
		}

		messages = append(messages, pkg.Message{
			Text:      reply,
			IsUser:    false,
			Timestamp: start,
			AudioSrc:  audioSrc,
		})
		c.WriteMessage(websocket.CloseMessage, []byte("close"))
	})

	r.POST("/chat", func(ctx *gin.Context) {
		var chatReq pkg.ChatRequest
		err := ctx.Bind(&chatReq)
		if err != nil {
			log.Println("[error]", err)
		}

		newMsg := pkg.Message{
			Text:      chatReq.Message,
			IsUser:    true,
			Timestamp: time.Now(),
		}
		messages = append(messages, newMsg)

		cfgAny, _ := ctx.Get("config")
		cfg, _ := cfgAny.(*config.Configuration)

		id := "answer-" + strconv.Itoa(int(time.Now().Unix()))
		if cfg.ChatModel == "ollama" {
			opts := api.NewChatOptions().WithTemperature(0.0).WithNumPredict(1000)
			tokenChan, err := api.ChatStreamOllama(
				*cfg.OllamaModel,
				messages,
				opts,
			)
			if err != nil {
				log.Println("[error] ChatStream", err)
			}
			answers[id] = tokenChan
		} else {
			tokenChan, err := api.ChatGPTChatStream(messages)
			if err != nil {
				log.Println("[error] ChatStream", err)
			}
			answers[id] = tokenChan
		}
		questions[id] = chatReq.Message

		templ.Handler(templates.ChatBubbleWithAnswer(newMsg, id)).ServeHTTP(ctx.Writer, ctx.Request)
	})
}
