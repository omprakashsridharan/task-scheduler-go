package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")
var validate *validator.Validate

type Config struct {
	Redis Redis `config:"redis"`
}

type Redis struct {
	Url string `config:"url" validate:"required"`
}

func InitConfig(envFilePath string, configFilePath string) (Config, error) {
	validate = validator.New()
	config := Config{}
	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Printf("Could not load environment file %s", err)
		}
	}

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return config, InvalidConfigFilePathError
	}

	if err := k.Load(file.Provider(configFilePath), json.Parser()); err != nil {
		log.Printf("error loading config: %v", err)
		return config, FileLoadError
	}

	err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "")), "_", ".", -1)
	}), nil)

	if err != nil {
		return config, FileLoadError
	}

	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "config"})
	if err != nil {
		return config, UnmarshalError
	}

	err = validate.Struct(config)
	if err != nil {
		log.Printf("Error while validating config: %s", err)
		return config, ValidationError
	}
	return config, nil
}
