package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// HelloWorld PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func HelloWorld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
