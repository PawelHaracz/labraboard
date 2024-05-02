package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"labraboard/internal/logger"
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
	l := logger.GetWitContext(g)
	FooBar(g)
	l.Info().Msg("Logged hello world")
}

func FooBar(ctx context.Context) {
	l := logger.GetWitContext(ctx)
	l.Info().Msg("called foobar")
}
