package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/BlackBart/server"
)

const (
	envKey  = "ENV"
	portKey = "PORT"
	local   = "local"
)

func main() {
	service, err := server.StartDefaultService()
	if err != nil {
		server.GetLogger().WithError(err).Panic("Can't initialize the service ...")
	}

	addRoutes(service)

	environment := os.Getenv(envKey)

	if environment == local {
		err = service.Run(":" + server.GetEnvPort(portKey))
	} else {
		err = nil
		service.RunAppEngine()
	}

	if err != nil {
		server.GetLogger().WithError(err).Panic("Can't start the server")
	}

}

func addRoutes(service *server.Service) {
	addPong(service)
}

func addPong(service *server.Service) {
	pingGroup := service.Group("/v1/blackbart/ping")
	{
		pingGroup.GET("", pong)
		pingGroup.POST("", pong)
		pingGroup.PUT("/:test", pong)
		pingGroup.PATCH("/:test", pong)
		pingGroup.DELETE("/:test", pong)
		pingGroup.GET("/:test", pong)
	}
}

func pong(c *gin.Context) {

	c.JSON(200, gin.H{
		"status":  "OK",
		"message": "pong",
	})
}

func pongFails(c *gin.Context) {
	response.SendInternalError(c, fmt.Errorf("a new error"))
}
