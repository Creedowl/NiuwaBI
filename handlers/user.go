package handlers

import (
	"fmt"
	"github.com/Creedowl/NiuwaBI/database/models"
	"github.com/gin-gonic/gin"
)

type RegisterParam struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// GetCurrentUser helper func
func GetCurrentUser(c *gin.Context) *models.User {
	user, exists := c.Get("id")
	if !exists {
		return nil
	}
	return user.(*models.User)
}

func Register(_ *gin.Context, param RegisterParam) (*models.User, error) {
	user, err := models.GetUserByName(param.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("user %s already existed", param.Username)
	}
	user, err = models.CreateUser(param.Username, param.Nickname, param.Password, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	return user, nil
}
