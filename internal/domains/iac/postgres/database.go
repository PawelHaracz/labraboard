package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	GormDB *gorm.DB
}

func NewDatabase(connectionString string) *Database {
	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	d, err := database.DB()
	if err != nil {
		panic("Failed to connect to database!")
	}
	if err = d.Ping(); err != nil {
		err := d.Close()
		if err != nil {
			panic(err)
		}
	}
	return &Database{GormDB: database}
}

func (db *Database) Close() error {
	sqlDB, err := db.GormDB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Close(); err != nil {
		return err
	}
	return nil
}
