package config

import (
	"os"
	"time"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// Init is startup configuration uploader and validator
func Init() {
	var config Config

	if err := config.Unmarshal(); err != nil {
		log.WithError(err).Fatal("Bad format in config file")
	}

	if err := config.Validate(); err != nil {
		log.WithError(err).Fatal("Bad configuration in config file")
	}

	if viper.GetViper().ConfigFileUsed() != "" {
		log.WithField("config", viper.GetViper().ConfigFileUsed()).Info("Complete load configuration")
	} else {
		log.Info("Complete load default configuration!")
	}
}

// Watch is realtime config watcher
func Watch() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			info, err := os.Stat(v.ConfigFileUsed())
			if os.IsNotExist(err) {
				log.WithError(err).Warn("Watcher lost configuration... Recreating!")
				if err := v.SafeWriteConfigAs(v.ConfigFileUsed()); err != nil {
					log.WithError(err).Warn("Can't create config dump. Disabling watcher!")
					return
				}
				continue
			} else if err != nil {
				log.WithError(err).Warn("Watcher work unhealthily... Disabling watcher!")
				return
			}

			if time.Unix(time.Now().Unix()-info.ModTime().Unix(), 0).Unix() < 5 {
				log.Info("Configuration change found! Updating running config...")
				var (
					config = Config{}
					lastv  = v.AllSettings()
				)
				if err := config.ReadIn(); err != nil {
					log.WithError(err).Warn("Can't read config file")
					continue
				}

				if err := config.Unmarshal(); err != nil {
					log.WithError(err).Warn("Bad format in config file")
					continue
				}

				if err := config.Validate(); err != nil {
					log.WithError(err).Warn("Bad configuration in config file")
					continue
				}

				if lastv["listen_port"].(int) != v.GetInt("listen_port") {
					v.Set("listen_port", lastv["listen_port"].(int))
					log.WithField("listen_port", v.GetInt("listen_port")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["db_name"] != v.GetString("database.db_name") {
					v.Set("database.db_name", lastv["database"].(map[string]interface{})["db_name"])
					log.WithField("db_name", v.GetString("database.db_name")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["host"] != v.GetString("database.host") {
					v.Set("database.host", lastv["database"].(map[string]interface{})["host"])
					log.WithField("host", v.GetString("database.host")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["password"] != v.GetString("database.password") {
					v.Set("database.password", lastv["database"].(map[string]interface{})["password"])
					log.WithField("password", "*******").Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["port"].(int) != v.GetInt("database.port") {
					v.Set("database.port", lastv["database"].(map[string]interface{})["port"])
					log.WithField("port", v.GetInt("database.port")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["sslmode"] != v.GetString("database.sslmode") {
					v.Set("database.sslmode", lastv["database"].(map[string]interface{})["sslmode"])
					log.WithField("sslmode", v.GetString("database.sslmode")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["user"] != v.GetString("database.user") {
					v.Set("database.user", lastv["database"].(map[string]interface{})["user"])
					log.WithField("user", v.GetString("database.user")).Warn("You can't change this value")
					continue
				}

				if lastv["database"].(map[string]interface{})["debug"] != v.GetBool("database.debug") {
					v.Set("database.debug", lastv["database"].(map[string]interface{})["debug"])
					log.WithField("debug", v.GetBool("database.debug")).Warn("You can't change this value")
					continue
				}

				if err := v.MergeInConfig(); err != nil {
					log.WithError(err).Warn("Can't merge configuration file with running config!")
				}
			}
		}
	}()
}
