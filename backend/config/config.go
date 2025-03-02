package config

import (
	"github.com/fox-one/pkg/db"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string `mapstructure:"port" default:"8000"`
	PprofAddr string `mapstructure:"pprof_addr" default:":28001"`

	// RedisConfig     *RedisConfig `mapstructure:"redis"`
	MixinConfig *MixinConfig `mapstructure:"mixin" required:"true"`
	// MongoConfig     *MongoConfig `mapstructure:"mongo"`
	DB *db.Config `mapstructure:"db" required:"true"`
}

type MixinConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	SessionID    string `mapstructure:"session_id"`
	PrivateKey   string `mapstructure:"private_key"`
	PinToken     string `mapstructure:"pin_token"`

	AppID             string `mapstructure:"app_id"`
	ServerPublicKey   string `mapstructure:"server_public_key"`
	SessionPrivateKey string `mapstructure:"session_private_key"`
	SpendKey          string `mapstructure:"spend_key"`

	EnableAutoReplay bool `mapstructure:"enable_auto_replay"`
}

func Init(filePath string, conf *Config) (err error) {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config file")
		return
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to unmarshal config")
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(_ fsnotify.Event) {
		log.Info().Msg("Config file changed")
		if err := viper.Unmarshal(conf); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal updated config")
		}
	})

	return nil
}
