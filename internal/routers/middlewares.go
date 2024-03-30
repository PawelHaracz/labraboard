package routers

import (
	"github.com/gin-gonic/gin"
	"labraboard/internal/helpers"
	"labraboard/internal/repositories"
)

func UnitedSetup(uow *repositories.UnitOfWork) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(helpers.UnitOfWorkSetup), uow)
		//c.Set("rc", rc)
		//c.Set("prefix", cfg.BucketPrefix)
	}
}
