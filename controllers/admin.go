package controllers

import (
	"gradio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddUserResponse struct {
	ID        string `json:"id"`
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
	Class     string `json:"class"`
	Password  string `json:"password,omitempty"`
	Rights    string `json:"rights"`
}

func AddUser(c *gin.Context) {
	var (
		db   = models.GetDB()
		data struct {
			GivenName string `json:"given_name" binding:"required"`
			Surname   string `json:"surname" binding:"required"`
			Class     string `json:"class" binding:"required"`
			Password  string `json:"password" binding:"omitempty"`
		}
		user models.User
	)

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.First(&user, "surname = ? AND class = ?", data.Surname, data.Class).RowsAffected != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this class and surname already exist"})
		return
	}

	user.GivenName = &data.GivenName
	user.Surname = data.Surname
	user.Class = data.Class
	if err := user.GenHash(data.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't create user hash from password"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't create user in database"})
		return
	}

	response := AddUserResponse{
		ID:        user.ID,
		GivenName: *user.GivenName,
		Surname:   user.Surname,
		Class:     user.Class,
		Password:  user.Password,
		Rights:    user.Rights,
	}

	c.JSON(http.StatusOK, gin.H{"user": response})
}

func DelStudent(c *gin.Context) {
	var (
		db   = models.GetDB()
		user models.User
		data struct {
			ID string `uri:"id" binding:"required,uuid"`
		}
	)

	if err := c.ShouldBindUri(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.First(&user, "id = ?", data.ID).RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this id not found"})
		return
	}

	if err := db.Model(&user).Association("Session").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on clear user session"})
		return
	}

	if err := db.Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on delete user"})
		return
	}

	// TODO: Добавить удаление контейнера сессии

	c.Status(http.StatusOK)
}

func CloseSession(c *gin.Context) {
	var (
		db   = models.GetDB()
		user models.User
		data struct {
			ID string `uri:"id" binding:"required,uuid"`
		}
	)

	if err := c.ShouldBindUri(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if db.First(&user, "id = ?", data.ID).RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this id not found"})
		return
	}

	if err := db.Model(&user).Association("Session").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on clear user session"})
		return
	}

	// TODO: Добавить удаление контейнера сессии

	c.Status(http.StatusOK)
}

func GetUsers(c *gin.Context) {
	var (
		db    = models.GetDB()
		users []models.User
	)

	db.Preload("Session").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}
