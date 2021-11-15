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

func Ping(_ *gin.Context, param PingResp) (*Pong, error) {
	return &Pong{
		Response: fmt.Sprintf("hello, %s", param.Name),
	}, nil
}

func Test(c *gin.Context) (*Pong, error) {
	u := GetCurrentUser(c)
	if u == nil || !u.IsAdmin() {
		return nil, errors.New("you are not admin")
	}
	return &Pong{
		Response: "hello, admin",
	}, nil
}
