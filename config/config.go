package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	ApiKey string
	ApiSecret string
	ApiPassphrase string
}

func init() {

}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateConfigFile() {
	filename := "config"
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	// Set default values as ""
	viper.SetDefault("ApiKey", "")
	viper.SetDefault("ApiSecret", "")
	viper.SetDefault("ApiPassphrase", "")

	if err := viper.WriteConfigAs(filename + ".yml"); err != nil {
		panic(err)
	}
}

func LoadConfig(filename string) (*Configuration, error) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("Error reading config file")
	}

	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.New("Unable to decode into struct")
	}

	return &config, nil
}

