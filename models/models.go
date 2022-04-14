package models

import (
	"fmt"
	"reflect"
	"time"

	gormlog "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Base contains common columns for all tables.
type Base struct {
	ID        string         `gorm:"type:uuid;primarykey;default:uuid_generate_v4()" json:"id" uri:"id" binding:"required,uuid"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// NewDBConnection creates a new database connection
func NewDBConnection() {
	var err error

	gormLogger := gormlog.New()
	gormLogger.SkipErrRecordNotFound = true
	gormLogger.SourceField = "line"
	gormLogger.SlowThreshold = time.Millisecond * 200

	gormConfig := gorm.Config{Logger: gormLogger}

	// Data Source Name for postgres connection
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.db_name"),
		viper.GetString("database.password"),
		viper.GetString("database.sslmode"),
	)

	log.WithField("dsn", dsn).Info("Connecting to database...")

	// Connecting to a database using ORM
	if db, err = gorm.Open(postgres.Open(dsn), &gormConfig); err != nil {
		log.WithError(err).Fatal("No connection to db")
	}

	log.Info("Database connected!")

	log.WithField("module", "uuid-ossp").Info("Add UUID extension to DB...")
	if err := db.Exec("create extension if not exists \"uuid-ossp\"").Error; err != nil {
		log.WithError(err).WithField("extension", "uuid-ossp").Fatal("Can't add extension")
	}

	models2Migrate := []interface{}{
		&User{},
		&Session{},
		&Grade{},
	}

	log.WithField("models", modelsNames(models2Migrate...)).Info("Migrating models...")
	if err := db.AutoMigrate(models2Migrate...); err != nil {
		log.WithError(err).Fatal("Can't migrate model to db")
	}
	log.Info("Models migrated...")

	var usersCount int64
	if db.Model(&User{}).Count(&usersCount); usersCount <= 0 {
		log.Info("Not found any users...")
		log.Info("Creating main admin user...")

		admin := &User{
			Surname: "admin",
			Rights:  "admin",
			Class:   "superclass",
		}

		if err := admin.GenHash("admin"); err != nil {
			log.WithError(err).Fatal("Can't generate default admin password")
		}

		if err := db.Create(&admin).Error; err != nil {
			log.WithError(err).Fatal("Unable to create administrator account")
		}
	}
}

func modelsNames(v ...interface{}) (names []string) {
	for _, model := range v {
		names = append(names, reflect.TypeOf(model).Elem().Name())
	}
	return
}

// GetDB returns the current open connection to the database
func GetDB() *gorm.DB {
	return db
}
