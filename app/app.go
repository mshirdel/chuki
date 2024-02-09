package app

import (
	"fmt"
	"github.com/mshirdel/chuki/config"
	"github.com/mshirdel/chuki/log"
)

type Application struct {
	configPath string
	Cfg        *config.Config
}

func New(configPath string) *Application {
	return &Application{
		configPath: configPath,
	}
}

func (a *Application) InitAll() error {
	if err := a.InitConfig(); err != nil {
		return err
	}
	a.InitLogger()

	return nil
}

func (a *Application) InitConfig() (err error) {
	if a.Cfg != nil {
		return
	}

	a.Cfg, err = config.InitViper(a.configPath)
	if err != nil {
		return fmt.Errorf("error in init config: %w", err)
	}

	return
}

func (a *Application) InitLogger() {
	cfg := a.Cfg.Logging
	log.InitLogrus(log.Config{Level: cfg.Level})
}
