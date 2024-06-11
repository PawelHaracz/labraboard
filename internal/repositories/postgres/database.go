package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"labraboard/internal/repositories/postgres/models"
)

type Database struct {
	GormDB *gorm.DB
}

func NewDatabase(connectionString string) *Database {
	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
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

func (db *Database) Migrate() {
	err := db.GormDB.AutoMigrate(
		&models.TerraformStateDb{},
		&models.IaCDb{},
		&models.IaCDeploymentDb{},
		&models.IaCPlanDb{})
	if err != nil {
		panic(errors.Wrap(err, "failed to migrate"))
	}
}
