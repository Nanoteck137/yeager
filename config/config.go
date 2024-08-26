package config

import (
	"os"

	"github.com/nanoteck137/yeager/core/log"
	"github.com/nanoteck137/yeager/types"
	"github.com/spf13/viper"
)

var AppName = "yeager"
var Version = "no-version"
var Commit = "no-commit"

type Config struct {
	ListenAddr      string `mapstructure:"listen_addr"`
	LibraryDir      string `mapstructure:"library_dir"`
}

func (c *Config) WorkDir() types.WorkDir {
	return types.WorkDir(c.LibraryDir)
}

func setDefaults() {
	viper.SetDefault("listen_addr", ":3000")
	viper.BindEnv("library_dir")
}

func validateConfig(config *Config) {
	hasError := false

	validate := func(expr bool, msg string) {
		if expr {
			log.Error("Config Validation", "err", msg)
			hasError = true
		}
	}

	// NOTE(patrik): Has default value, here for completeness
	validate(config.ListenAddr == "", "listen_addr needs to be set")
	validate(config.LibraryDir == "", "library_dir needs to be set")

	if hasError {
		log.Fatal("Config not valid")
	}
}

var ConfigFile string
var LoadedConfig Config

func InitConfig() {
	setDefaults()

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(AppName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("Failed to load config", "err", err)
	}

	err = viper.Unmarshal(&LoadedConfig)
	if err != nil {
		log.Error("Failed to unmarshal config: ", err)
		os.Exit(-1)
	}

	log.Debug("Current Config", "config", LoadedConfig)
	validateConfig(&LoadedConfig)
}
