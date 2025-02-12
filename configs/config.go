package configs

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
)

//go:embed config.yaml
var Configurations []byte

type App struct {
	Token     string
	Debug     bool
	ChannelID int64
}
type Config struct {
	App *App
}

func NewConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(Configurations)); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
