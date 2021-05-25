package database

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	DB() interface{}
	Migrate() error
}

type dbGORM struct {
	db *gorm.DB
}

func NewDBGORM(uri string, cfg *gorm.Config) Database {
	db, err := gorm.Open(postgres.Open(uri), cfg)
	if err != nil {
		panic("failed to connect to database")
	}
	return dbGORM{db}
}

func (g dbGORM) Migrate() error {
	if err := g.db.AutoMigrate(entity.User{}); err != nil {
		return err
	}
	if err := g.db.AutoMigrate(entity.Credentials{}); err != nil {
		return err
	}
	if err := g.db.AutoMigrate(entity.PublicInfo{}); err != nil {
		return err
	}
	if err := g.db.AutoMigrate(entity.Video{}); err != nil {
		return err
	}
	if err := g.db.AutoMigrate(entity.Provider{}); err != nil {
		return err
	}
	return nil
}

func (g dbGORM) DB() interface{} {
	return g.db
}
