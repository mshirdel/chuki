package app

import (
	"fmt"
	"github.com/mshirdel/chuki/app/db"
	"github.com/mshirdel/chuki/config"
	"github.com/mshirdel/chuki/log"
)

type Application struct {
	configPath string
	Cfg        *config.Config
	DB         *db.DB
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
	a.ShowBanner()
	a.InitLogger()
	if err := a.InitDatabase(); err != nil {
		return err
	}

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

func (a *Application) InitDatabase() error {
	if a.DB != nil {
		return nil
	}

	a.DB = db.New(&a.Cfg.Database)
	err := a.DB.InitAll()
	if err != nil {
		return fmt.Errorf("error in init database: %w", err)
	}

	return nil
}

func (a *Application) ShowBanner() {
	fmt.Println(`
      _           _    _ 
  ___| |__  _   _| | _(_)
 / __| '_ \| | | | |/ / |
| (__| | | | |_| |   <| |
 \___|_| |_|\__,_|_|\_\_|`)
}
