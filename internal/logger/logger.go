package logger

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/net/context"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var once sync.Once
var hasInitialized = false
var log zerolog.Logger

const (
	CorrelationIdKey = "correlation_id"
	RequestIdKey     = "request_id"
)

func Init(logLevel int8, prettyLogs bool) {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

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

		log = zerolog.New(os.Stdout).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()

		if prettyLogs {
			var output io.Writer = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
			log = log.Output(output)
		}

		hasInitialized = true
	})
}

func Get() zerolog.Logger {
	if !hasInitialized {
		panic(errors.New("Logger is not initialized"))
	}
	return log
}

func GetWitContext(ctx context.Context) zerolog.Logger {
	l := Get()

	correlationId := ""
	requestId := ""
	value := ctx.Value(CorrelationIdKey)
	if value != nil {
		correlationId = value.(string)
	}
	value = ctx.Value(RequestIdKey)
	if value != nil {
		requestId = value.(string)
	}
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		if correlationId != "" {
			c = c.Str(CorrelationIdKey, correlationId)
		}

		if requestId != "" {
			c = c.Str(RequestIdKey, requestId)
		}
		return c
	})

	return l
}
