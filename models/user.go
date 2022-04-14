package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is table of shop users
type User struct {
	Base
	Rights    string   `json:"rights" gorm:"default:student"`
	Surname   string   `json:"surname" gorm:"size:128;not null" `
	Class     string   `json:"class" gorm:"not null"`
	GivenName *string  `json:"given_name,omitempty" gorm:"size:128"`
	Session   *Session `json:"session,omitempty"`
	Grades    []Grade  `json:"grades,omitempty"`
	Hash      string   `json:"-" gorm:"not null"`
	Password  string   `json:"password,omitempty" gorm:"-"`
}

func (u *User) Get(id string) error {
	if db.First(&u, "id = ?", id).RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// GenHash is generate password hash to this model
func (u *User) GenHash(pass string) (err error) {
	if pass == "" {
		rand.Seed(time.Now().UnixNano())
		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
		length := 8
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		u.Password = b.String()
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost); err == nil {
		u.Hash = string(hash)
	}

	return err
}

type Session struct {
	Base
	UserID        string `json:"-"`
	Port          uint   `json:"-"`
	ContainerID   string
	ConnectionURL string `json:"connection_url"`
}

type Grade struct {
	UserID string
	Mark   int
}
