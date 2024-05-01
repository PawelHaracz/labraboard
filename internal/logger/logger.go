package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var once sync.Once
var hasInitialized = false
var log zerolog.Logger

var contextLogger string = "GIN_LOGGER"

func Init(logLevel int8) {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()

		hasInitialized = true
	})
}

func Get() zerolog.Logger {
	if !hasInitialized {
		panic(errors.New("Logger is not initialized"))
	}
	return log
}

func GetGinLogger(c *gin.Context) zerolog.Logger {
	correlationId := TakeValueFromHeader(c, "X-Correlation-ID")
	requestId := TakeValueFromHeader(c, "X-Request-ID")
	l := Get()
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.
			Str("correlation_id", correlationId).
			Str("X-Request-ID", requestId)
	})

	return l
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		l := GetGinLogger(c)
		path := c.Request.URL.Path
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

		c.Set(contextLogger, l)
	}
}

func TakeValueFromHeader(c *gin.Context, header string) string {
	headerValues, ok := c.Request.Header[header]
	var value = ""
	if !ok {
		value = xid.New().String()
		c.Request.Header.Add(header, value)
	} else {
		value = headerValues[0]
	}

	return value
}
