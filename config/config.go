package config

import (
	"bytes"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Logging Logging `yaml:"logging"`
}

type Logging struct {
	Level string `yaml:"level"`
}

const envPrefix = "chuki"

func InitViper(configPath string) (*Config, error) {
	var c Config

	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(builtinConfig))); err != nil {
		return nil, fmt.Errorf("error loading default config: %w", err)
	}

	v.SetConfigFile(configPath)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	switch err := v.MergeInConfig(); err.(type) {
	case nil:
	case *os.PathError:
		logrus.Infof("config file (%s) not found, using default and env variables", configPath)
	default:
		logrus.Warnf("failed to load config file: %s", err)
	}

	if err := v.UnmarshalExact(&c); err != nil {
		return nil, fmt.Errorf("faild to unmarshal config into struct: %w", err)
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("faild to validate config sructure: %w", err)
	}

	return &c, nil
}

func (c *Config) Validate() error {
	return validator.New().Struct(c)
}
