package client

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	api "github.com/orov-io/BlackBeard"
)

const (
	portkey      = "PORT"
	serviceKey   = "SERVICE_BASE_PATH"
	v1           = "v1"
	pingEndpoint = "/ping"
)

var service = os.Getenv(serviceKey)

// Ping make a call to the is_alive endpoint.
func Ping(c *gin.Context) (*PongResponse, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithVersion(v1).
		ToService(service)
	resp, err := client.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := new(PongResponse)
	err = response.ParseTo(resp, &pong)

	return pong, err
}

// PongResponse is the expected response to ping request
type PongResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
