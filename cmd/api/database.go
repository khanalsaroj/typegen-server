package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/khanalsaroj/typegen-server/internal/config"
	"github.com/khanalsaroj/typegen-server/internal/domain"

	gormlogger "gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func connectDatabase(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dbPath := cfg.Database.Filepath

	logLevel := gormlogger.Silent
	if cfg.App.Environment == "dev" {
		logLevel = gormlogger.Info
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil && filepath.Dir(dbPath) != "." {
		log.Fatal("err")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	log.Info("SQLite connected successfully")
	return db, nil
}

func runMigrations(db *gorm.DB, log *zap.Logger) error {
	log.Info("Running database migrations (DTO-based)...")

	if err := db.AutoMigrate(
		&domain.DatabaseConnection{},
	); err != nil {
		return err
	}

	log.Info("DTO migrations completed successfully")
	return nil
}
