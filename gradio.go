package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"gradio/config"
	"gradio/controllers"
	"gradio/middleware"
	"gradio/models"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting gradio...")

	config.Init()
	config.Watch()

	r := gin.Default()
	r.Use(middleware.AllowCORSConfig())
	models.NewDBConnection()
	pullGnuImage()

	r.GET("ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "pong"}) })

	// Роуты сессий студентов
	session := r.Group("session")
	{
		session.POST("", controllers.GenerateSession)
		session.GET(":id", controllers.GetStatusOfSession)
		session.DELETE(":id", controllers.StopAndDeleteSession)
	}

	admin := r.Group("admin")
	{
		// Управление студентами
		users := admin.Group("users")
		users.GET("", controllers.GetUsers)
		users.GET(":id", controllers.NotImplemented)
		users.POST("", controllers.AddUser)
		users.PUT(":id", controllers.NotImplemented)
		users.DELETE(":id", controllers.DelStudent)
		// Управление оценками студентов
		users.GET(":id/grades", controllers.NotImplemented)
		users.POST(":id/grades", controllers.NotImplemented)
		users.PUT(":id/grades/:grade_id", controllers.NotImplemented)
		users.DELETE(":id/grades/:grade_id", controllers.NotImplemented)
		// Управление сессиями студентов
		users.POST(":id/session", controllers.NotImplemented)
		users.DELETE(":id/session", controllers.CloseSession)
	}

	if _, err := net.Dial("tcp", "localhost:"+viper.GetString("listen_port")); err == nil {
		log.WithError(err).WithField("port", viper.GetString("listen_port")).Fatal("Can't bind port")
	}

	log.WithField("port", viper.GetString("listen_port")).Info("Starting server...")
	if err := r.Run(":" + viper.GetString("listen_port")); err != nil {
		log.WithError(err).Fatal("AAAA Panic, Server 1$ D0wn.... jco8*")
	}
}

func pullGnuImage() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	authConfig := types.AuthConfig{
		Username: viper.GetString("registry.user"),
		Password: viper.GetString("registry.password"),
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, viper.GetString("registry.image"), types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}
