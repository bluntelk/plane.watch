package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
)

type BotConfig struct {
	Token          string
	HereMapsApiKey string
}

func Load() *BotConfig {
	conf := BotConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("Token", "your-bot-token-here")
	viper.AutomaticEnv()

	binaryPath, err := os.Executable()
	if nil != err {
		log.Fatalln(err)
	}
	binaryDir := path.Dir(binaryPath)
	if !strings.Contains(binaryPath, "go-build") {
		viper.AddConfigPath(binaryDir)
	}
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// write our default config
			if err := viper.SafeWriteConfig(); nil != err {
				log.Fatalln(err)
			}
			log.Fatalf("No config found, example file written to: %s", binaryDir)
		} else {
			log.Fatalln(err)
		}
	}

	conf.Token = viper.GetString("Token")
	conf.HereMapsApiKey = viper.GetString("HereMapsApiKey")

	return &conf
}
