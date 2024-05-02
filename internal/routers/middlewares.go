package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"labraboard/internal/helpers"
	"labraboard/internal/logger"
	"labraboard/internal/repositories"
	"time"
)

const (
	CorrelationId string = "X-Correlation-ID"
)

func UnitedSetup(uow *repositories.UnitOfWork) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(helpers.UnitOfWorkSetup), uow)
		//c.Set("rc", rc)
		//c.Set("prefix", cfg.BucketPrefix)
	}
}

func UseCorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValues, ok := c.Request.Header[CorrelationId]
		var correlationId = ""
		if !ok {
			correlationId = xid.New().String()
			c.Request.Header.Add(CorrelationId, correlationId)
		} else {
			correlationId = headerValues[0]
		}
		c.Header(CorrelationId, correlationId)
		c.Next()
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		l := logger.GetWitContext(c)
		path := c.Request.URL.Path
		correlationId := TakeValueFromHeader(c, "X-Correlation-Id")
		requestId := TakeValueFromHeader(c, "X-Request-Id")
		c.Set(logger.CorrelationIdKey, correlationId)
		c.Set(logger.RequestIdKey, requestId)
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			if correlationId != "" {
				c = c.Str(logger.CorrelationIdKey, correlationId)
			}

			if requestId != "" {
				c = c.Str(logger.RequestIdKey, requestId)
			}
			return c

		})
		c.Header("X-Correlation-Id", correlationId)
		c.Next()

		// access the status we are sending
		status := c.Writer.Status()

		l.
			Info().
			Str("path", path).
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.RequestURI()).
			Str("user_agent", c.Request.UserAgent()).
			Dur("elapsed_ms", time.Since(t)).
			Str("client_ip", c.ClientIP()).
			Int("status_code", status).
			Msg("incoming request")
	}
}

func TakeValueFromHeader(c *gin.Context, header string) string {
	headerValues, ok := c.Request.Header[header]
	if ok {
		return headerValues[0]
	}
	return ""
}
