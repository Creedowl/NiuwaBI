package main

import (
	"github.com/Creedowl/NiuwaBI/database/models"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "niuwa",
		Key:         []byte(utils.Cfg.Auth.JwtKey),
		Timeout:     time.Second * time.Duration(utils.Cfg.Auth.Timeout),
		MaxRefresh:  time.Second * time.Duration(utils.Cfg.Auth.MaxRefresh),
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":          v.ID,
					"name":        v.Name,
					"nickname":    v.Nickname,
					"permissions": v.Permissions,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return &models.User{
				BaseModel:   models.BaseModel{ID: uint(claims["id"].(float64))},
				Name:        claims["name"].(string),
				Permissions: uint(claims["permissions"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user, err := models.GetAuthUser(username, password)
			if err != nil {
				logrus.Errorf("failed to get auth user: %v", err)
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
				"data":    nil,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":         http.StatusOK,
				"access_token": token,
				"expire":       expire.Format(time.RFC3339),
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
