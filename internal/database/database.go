package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/watzon/0x45/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func New(config *config.Config, gormConfig *gorm.Config) (*Database, error) {
	var dialector gorm.Dialector

	switch config.Database.Driver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s port=%d dbname=%s sslmode=%s",
			config.Database.Host,
			config.Database.User,
			config.Database.Password,
			config.Database.Port,
			config.Database.Name,
			config.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(config.Database.Name)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Database.Driver)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{db}, nil
}

func (d *Database) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *Database) Migrate(config *config.Config) error {
	return RunMigrations(d.DB)
}
