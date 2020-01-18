package client

import (
	"github.com/orov-io/BlackBeard"
	"github.com/orov-io/BlackBart/response"
	"github.com/gin-gonic/gin"
)

const (
	portkey      = "ADMIN_PORT"
	service      = "admin"
	lastVersion  = "v1"
	flagEndpoint = "/flags/user"
	userEndpoint = "user"
	groupSubfix  = "group"
	pingEndpoint = "/ping"
)

// Ping make a call to the is_alive endpoint.
func Ping(c *gin.Context) (*PongResponse, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithVersion(lastVersion).
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
