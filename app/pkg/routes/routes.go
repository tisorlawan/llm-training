package routes

import (
	"app/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	Phi        = "phi:latest"
	Llama      = "llama3"
	NeuralChat = "neural-chat:latest"
)

var upgrader = websocket.Upgrader{}

func CreateApp(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./assets/dist")
	r.Static("/images", "./images")
	r.Static("/audios", "./audios")

	r.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
	})

	RoutesChat(r, "")
	RoutesSummary(r, "")

	return r
}
