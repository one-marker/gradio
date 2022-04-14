package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var v = viper.GetViper()

// Init entrypoint of configuration
func init() {
	v.SetDefault("database.host", "db")
	v.SetDefault("database.db_name", "gradio")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.debug", true)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("listen_port", 3000)
	v.SetDefault("external_schema", "http")
	v.SetDefault("external_host", "localhost")
	v.SetDefault("registry.image", "gosgradio/gradio")

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		DisableQuote:    true,
		TimestampFormat: time.Stamp,
	})

	v.AutomaticEnv()
	v.SetConfigName("gradio")
	v.AddConfigPath("/etc/gradio/")
	v.AddConfigPath("$HOME/.gradio/")
	v.AddConfigPath(".")

	if err := (&Config{}).ReadIn(); err != nil {
		log.WithError(err).Fatal("Can't init config file")
	}
}

// Config is a structure off all settings in BirkaAPI, that contains validator for checks
type Config struct {
	ListenPort int `mapstructure:"listen_port" validate:"required,numeric,gte=1,lte=65535"`
	Database   struct {
		Host        string `mapstructure:"host" validate:"required,ip|hostname"`
		Port        int    `mapstructure:"port" validate:"required,numeric,gte=1,lte=65535"`
		User        string `mapstructure:"user" validate:"required"`
		DBName      string `mapstructure:"db_name" validate:"required"`
		Password    string `mapstructure:"password" validate:"required"`
		SSLMode     string `mapstructure:"sslmode" validate:"required,oneof=enable disable verify-full"`
		SSLRootCert string `mapstructure:"sslrootcert" validate:"omitempty"`
		Debug       bool   `mapstructure:"debug"`
	} `mapstructure:"database" validate:"required,dive"`
	Registry struct {
		Image    string `mapstructure:"image" validate:"required"`
		User     string `mapstructure:"user" validate:"omitempty"`
		Password string `mapstructure:"password" validate:"omitempty"`
	} `mapstructure:"registry" validate:"required,dive"`
	ExternalHost   string `mapstructure:"external_host" validate:"required,hostname"`
	ExternalSchema string `mapstructure:"external_schema" validate:"required,oneof=http https"`
}

// Validate base check config variables
func (c *Config) Validate() error {
	return validator.New().Struct(c)
}

// Unmarshal unmarshal config
func (c *Config) Unmarshal() error {
	return v.UnmarshalExact(c)
}

// ReadIn is reading configuration into file
func (c *Config) ReadIn() (err error) {
	err = v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		if err = v.SafeWriteConfigAs("gradio.yml"); err != nil {
			return
		}
		return
	}
	return
}
