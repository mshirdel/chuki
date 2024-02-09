package db

import (
	"fmt"
	"github.com/mshirdel/chuki/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewOrCreate(cfg *config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       cfg.DSN(),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: cfg.SkipInitializeWithVersion,
	}), &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold:             cfg.Logger.SlowThreshold,
			LogLevel:                  cfg.Logger.GormLogLevel(),
			Colorful:                  cfg.Logger.Colorful,
			IgnoreRecordNotFoundError: cfg.Logger.IgnoreRecordNotFoundError,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("can't connect to database %s", cfg.DSN())
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("can't get sql database %s : %w", cfg.DSN(), err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)

	var connectionID int
	tx := db.Raw("SELECT CONNECTION_ID()").Scan(&connectionID)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, fmt.Errorf("can't ping database %s: %w", cfg.DSN(), err)
	}

	logrus.Debugf("[PING] connected to MySQL database with connection id: %d", connectionID)

	return db, nil
}
