package controllers

import (
	"context"
	"fmt"
	"gradio/models"
	"gradio/tools"
	"net"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetStatusOfSession(c *gin.Context) {
	var (
		db   = models.GetDB()
		data struct {
			ID string `uri:"id" binding:"required,uuid"`
		}
		session models.Session
	)

	if err := c.ShouldBindUri(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.First(&session, "id = ?", data.ID).RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "session with this id not found"})
		return
	}

	if _, err := net.Dial("tcp", viper.GetString("external_host")+":"+strconv.Itoa(int(session.Port))); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "offline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "online",
		"connection_url": session.ConnectionURL,
		"port":           session.Port,
	})
}

func GenerateSession(c *gin.Context) {
	var (
		db   = models.GetDB()
		data struct {
			Surname string `json:"surname" binding:"required"`
			Class   string `json:"class" binding:"required"`
		}
		user models.User
	)

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.Preload("Session").First(&user, "surname = ? AND class = ?", data.Surname, data.Class).RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this class and surname not found"})
		return
	}

	if user.Session == nil {
		availablePort := tools.GetEmptyPort()
		containerID := runGnuContainer(strconv.Itoa(availablePort))

		user.Session = &models.Session{
			Port:          uint(availablePort),
			ContainerID:   containerID,
			ConnectionURL: fmt.Sprintf("vnc://vuc@%s:%d", viper.GetString("external_host"), availablePort),
		}
		db.Save(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"connection_url": user.Session.ConnectionURL,
		"surname":        user.Surname,
		"class":          user.Class,
		"session_id":     user.Session.ID,
	})
}

func runGnuContainer(port string) (containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	exposedPorts, portBindings, err := nat.ParsePortSpecs([]string{port + ":5900"})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        viper.GetString("registry.image"),
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func StopAndDeleteSession(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
