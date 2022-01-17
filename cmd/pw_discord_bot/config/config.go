package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

type BotConfig struct {
	Token          string
	HereMapsApiKey string
}

func Load(configFile string) *BotConfig {
	conf := BotConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("Token", "your-bot-token-here")
	viper.AutomaticEnv()

	binaryDir := "."
	if "" == configFile {
		binaryPath, err := os.Executable()
		if nil != err {
			log.Error().Err(err).Msg("Cannot determine executables path")
			return nil
		}
		binaryDir = path.Dir(binaryPath)
		if !strings.Contains(binaryPath, "go-build") {
			viper.AddConfigPath(binaryDir)
		}
		viper.AddConfigPath(".")
	} else {
		log.Debug().Str("file", configFile).Msg("Using specified config file")
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// write our default config
			if err := viper.SafeWriteConfig(); nil != err {
				log.Error().Err(err).Msg("Failed to write config")
				return nil
			}
			log.Info().Msgf("No config found, example file written to: %s", binaryDir)
		} else {
			log.Error().Err(err).Msg("")
			return nil
		}
	}

	conf.Token = viper.GetString("Token")
	conf.HereMapsApiKey = viper.GetString("HereMapsApiKey")

	return &conf
}
