package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Options struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (o Options) String() string {
	return fmt.Sprintf("MaxOpenConns: %d MaxIdleConns: %d ConnMaxLifetime: %v", o.MaxOpenConns, o.MaxIdleConns, o.ConnMaxLifetime)
}

func Connect(dsn string, opts *Options) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil && opts != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		log.Printf("Setting up connection pool: %v", opts)
		sqlDB.SetMaxOpenConns(opts.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(opts.ConnMaxLifetime)
		sqlDB.SetMaxIdleConns(opts.MaxIdleConns)
	}
	return db, err
}
