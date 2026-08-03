package database

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.DbConfig) (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort, cfg.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	if tx := db.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS idx_role_user 
    ON welcome_roles (role_id, COALESCE(user_id, 'NULL_REPLACEMENT')) 
    WHERE deleted_at IS NULL
	`); tx.Error != nil {
		return tx.Error
	}

	return nil
}
