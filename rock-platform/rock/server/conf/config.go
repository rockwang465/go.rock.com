package conf

import "github.com/spf13/viper"

type Config struct {
	Viper *viper.Viper
}

func GetConfig() *Config {
	return &Config{
		Viper: viper.GetViper(),
	}
}
