package routers

import (
	"github.com/gin-gonic/gin"
	"labraboard/internal/helpers"
	"labraboard/internal/repositories"
	"labraboard/internal/repositories/postgres"
)

func UnitedSetup(db *postgres.Database) gin.HandlerFunc {
	uow, err := repositories.NewUnitOfWork(
		repositories.WithIaCRepositoryDbRepository(db),
		repositories.WithTerraformStateDbRepository(db))
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Set(string(helpers.UnitOfWorkSetup), uow)
		//c.Set("rc", rc)
		//c.Set("prefix", cfg.BucketPrefix)
	}
}
