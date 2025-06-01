package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	output.FormatLevel = func(i interface{}) string {
		switch i {
		case "info":
			return "\033[32mINFO\033[0m" // Verde
		case "error":
			return "\033[31mERROR\033[0m" // Vermelho
		case "debug":
			return "\033[33mDEBUG\033[0m" // Amarelo
		default:
			return i.(string)
		}
	}

	Logger = zerolog.New(output).With().Timestamp().Logger()
	log.Logger = Logger
}