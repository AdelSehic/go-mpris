package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	once sync.Once
)

// InitLogger initializes the global logger with the given log level and optional writer.
// If no writer is provided, it defaults to pretty console output to stderr.
func InitLogger(level zerolog.Level, writer ...io.Writer) {
	once.Do(func() {
		zerolog.SetGlobalLevel(level)

		var out io.Writer
		if len(writer) > 0 && writer[0] != nil {
			out = writer[0]
		} else {
			out = zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			}
		}

		log.Logger = zerolog.New(out).With().Timestamp().Logger()
		log.Logger = log.Output(out) // <- key for pretty output
	})
}

// GetLogger returns the global zerolog logger.
func GetLogger() zerolog.Logger {
	return log.Logger
}
