package config

import (
	"bytes"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Logging  Logging  `yaml:"logging"`
	Database Database `yaml:"database"`
}

type Database struct {
	Host                      string        `yaml:"host"`
	Port                      int           `yaml:"port"`
	User                      string        `yaml:"user"`
	Password                  string        `yaml:"password"`
	DBName                    string        `yaml:"dbname"`
	Charset                   string        `yaml:"charset"`
	Collation                 string        `yaml:"collation"`
	ParseTime                 bool          `yaml:"parse_time"`
	Location                  string        `yaml:"location"`
	MaxLifeTime               time.Duration `yaml:"max_life_time"`
	MaxIdleTime               time.Duration `yaml:"max_idle_time"`
	MaxOpenConnections        int           `yaml:"max_open_connections"`
	MaxIdleConnections        int           `yaml:"max_idle_connections"`
	SkipInitializeWithVersion bool          `yaml:"skip_initialize_with_version"`
	Logger                    DBLogger      `yaml:"logger"`
}

func (d *Database) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=%t&charset=%s&collation=%s&loc=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.ParseTime, d.Charset, d.Collation, url.PathEscape(d.Location))
}

type DBLogger struct {
	SlowThreshold             time.Duration `yaml:"slow_threshold"`
	Level                     string        `yaml:"level"`
	Colorful                  bool          `yaml:"colorful"`
	IgnoreRecordNotFoundError bool          `yaml:"ignore_record_not_found_error"`
}

func (l DBLogger) GormLogLevel() logger.LogLevel {
	switch l.Level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

type Logging struct {
	Level string `yaml:"level"`
}

const envPrefix = "chuki"

func InitViper(configPath string) (*Config, error) {
	var c Config

	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader([]byte(builtinConfig))); err != nil {
		return nil, fmt.Errorf("error loading default config: %w", err)
	}

	if configPath != "" {
		v.SetConfigFile(configPath)
		switch err := v.MergeInConfig(); err.(type) {
		case nil:
		case *os.PathError:
			logrus.Infof("config file (%s) not found, using default and env variables", configPath)
		default:
			logrus.Warnf("failed to load config file: %s", err)
		}
	}

	err := v.Unmarshal(&c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	})
	if err != nil {
		logrus.WithError(err).Errorf("error parsing configs to Config struct")
	}

	return &c, nil
}

func (c *Config) Validate() error {
	return validator.New().Struct(c)
}
