package database

import (
	"fmt"
	"grf/core/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func dialector(config *config.Config) gorm.Dialector {
	if config.DBVendor == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
		return mysql.Open(dsn)
	} else if config.DBVendor == "sqlite" {
		return sqlite.Open(config.DBName)
	} else {
		log.Fatalf("Unsupported database vendor: %s", config.DBVendor)
		return nil
	}
}

func gormConfig(config *config.Config) *gorm.Config {
	var logger = gormLogger.Warn
	if config.DBLogLevel == "info" {
		logger = gormLogger.Info
	} else if config.DBLogLevel == "error" {
		logger = gormLogger.Error
	} else if config.DBLogLevel == "silent" {
		logger = gormLogger.Silent
	}
	return &gorm.Config{
		Logger: gormLogger.Default.LogMode(logger),
	}
}

func PerformMigration(db *gorm.DB, config *config.Config, dst ...interface{}) error {
	if config.DBMigrate {
		log.Println("Running AutoMigrate ")

		err := db.AutoMigrate(dst...)

		if err != nil {
			return fmt.Errorf("failed to perform migrations: %w", err)
		}
	}
	return nil
}

func ConnectDB(config *config.Config) (*gorm.DB, error) {

	log.Printf("Connecting to database %s:%s/%s", config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(dialector(config), gormConfig(config))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.DBMaxIdle)
	sqlDB.SetMaxOpenConns(config.DBMaxOpened)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.DBMaxIdle) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DBMaxLifeTimeSeconds) * time.Second)

	log.Println("Database connected successfully")
	return db, nil
}
