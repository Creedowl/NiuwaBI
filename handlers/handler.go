package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type PingResp struct {
	Name string `json:"name" form:"name"`
}

type Pong struct {
	Response string `json:"response"`
}

func Ping(c *gin.Context, param PingResp) (*Pong, error) {
	return &Pong{
		Response: fmt.Sprintf("hello, %s", param.Name),
	}, nil
}

func Test(c *gin.Context) (*Pong, error) {
	return nil, errors.New("error test")
}
