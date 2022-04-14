package main

import (
	"strings"
	"time"

	"gradio/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// JWT is used to pass through jwt
var JWT = func() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(viper.GetString("JWT_SECRET")),
		MaxRefresh: 720 * time.Hour,
		Timeout:    time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(models.User); ok {
				return jwt.MapClaims{
					"id": u.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			user := models.User{}
			user.ID = claims["id"].(string)
			return &user
		},
		Authenticator: authenticate,
		Authorizator:  authorizate,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.Status(code)
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
	})
	if err != nil {
		log.WithError(err).Fatal("JWT Error!")
	}

	return authMiddleware
}()

// authorizate user authorization handler
func authorizate(data interface{}, c *gin.Context) bool {
	user := data.(*models.User)
	if err := user.Get(user.ID); err != nil {
		log.WithError(err).Warn("Can't authorize user")
		return false
	}
	return true
}

// authenticate user authentication handler
func authenticate(c *gin.Context) (interface{}, error) {
	var (
		db       = models.GetDB()
		user     models.User
		authData struct {
			Login    string `json:"login" binding:"required,email"`
			Password string `json:"password" binding:"required,upwd"`
		}
	)

	if err := c.ShouldBindJSON(&authData); err != nil {
		return "", err
	}

	lowerLogin := strings.ToLower(authData.Login)
	if db.First(&user, "Surname = ?", lowerLogin).RowsAffected == 0 {
		return "", jwt.ErrFailedAuthentication
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(authData.Password))
	if err != nil {
		return "", jwt.ErrFailedAuthentication
	}

	return user, nil
}
