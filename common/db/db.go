package db

import (
	"fmt"

	"github.com/MonsieurTa/hypertube/common/infrastructure/database"
	"gorm.io/gorm"
)

type PSQLConfig struct {
	Host     string
	User     string
	Password string
	Db       string
	Port     string
}

func InitDB(psqlCfg *PSQLConfig, gormCfg *gorm.Config) database.Database {
	format := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"

	db_uri := fmt.Sprintf(format,
		psqlCfg.Host,
		psqlCfg.User,
		psqlCfg.Password,
		psqlCfg.Db,
		psqlCfg.Port)

	db := database.NewDBGORM(db_uri, &gorm.Config{})

	err := db.Migrate()
	if err != nil {
		panic("failed to migrate")
	}
	return db
}
