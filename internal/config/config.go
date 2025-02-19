// @description:
// @file: config.go
// @date: 2021/11/16

// Package config 读取配置文件。
package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/saltfishpr/redis-viewer/internal/util"

	"github.com/spf13/viper"
)

// Config represents the main config for the application.
type Config struct {
	Addrs    []string `mapstructure:"addrs"`
	DB       int      `mapstructure:"db"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`

	// sentinel
	MasterName string `mapstructure:"master_name"`

	Count int64 `mapstructure:"count"` // default 20
}

// LoadConfig loads a users config and creates the config if it does not exist.
func LoadConfig() {
	configPath, err := util.GetHomeDirectory()
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS != "windows" {
		configPath = filepath.Join(configPath, ".config", "redis-viewer")
		if err = util.CreateDirectory(configPath); err != nil {
			log.Fatal(err)
		}
	}

	viper.AddConfigPath(configPath)

	viper.SetConfigName("redis-viewer")
	viper.SetConfigType("yml")

	viper.SetDefault("mode", "client")
	viper.SetDefault("count", 20)

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfig()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
