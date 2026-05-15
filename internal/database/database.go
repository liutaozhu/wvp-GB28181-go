package database

import (
	"fmt"
	"time"

	"wvp-pro-go/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig, log *zap.Logger) error {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormLogger := logger.New(
		&zapWriter{log: log},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
	log.Info("database initialized", zap.String("driver", cfg.Driver))
	return nil
}

// AutoMigrate runs auto migration for given models
func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}

// zapWriter implements gorm logger.Writer interface
type zapWriter struct {
	log *zap.Logger
}

func (w *zapWriter) Printf(format string, args ...interface{}) {
	w.log.Info(fmt.Sprintf(format, args...))
}
