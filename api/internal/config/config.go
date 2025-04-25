package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func MustInit(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	fmt.Println(fmt.Sprintf("CONFIG_PATH=%s", configPath))
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
