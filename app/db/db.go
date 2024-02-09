package db

import (
	"github.com/mshirdel/chuki/config"
	"gorm.io/gorm"
)

type DB struct {
	cfg   *config.Database
	Chuki *gorm.DB
}

func New(cfg *config.Database) *DB {
	return &DB{cfg: cfg}
}

func (d *DB) InitAll() error {
	if err := d.initChukiDB(); err != nil {
		return err
	}

	return nil
}

func (d *DB) initChukiDB() error {
	var err error
	if d.Chuki != nil {
		return err
	}

	d.Chuki, err = NewOrCreate(d.cfg)

	return err
}
