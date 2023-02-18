package main

import (
	"flag"
	"os"

	"github.com/itchyny/gojq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultLogLevel = zerolog.DebugLevel

func main() {
	InitLogging(defaultLogLevel)

	hp := GetHandlerParams()

	RunServer(&hp)
}

func GetHandlerParams() HandlerParams {
	hp := HandlerParams{
		expression: os.Getenv("JQ_EXPRESSION"),
		host:       os.Getenv("HOST"),
		port:       os.Getenv("PORT"),
		scheme:     "", // TODO: https
	}

	flag.StringVar(&hp.expression, "jq", hp.expression, "The jq filter")
	flag.StringVar(&hp.host, "a", hp.host, "The server's host configuration, default \"\"")
	flag.StringVar(&hp.port, "p", hp.port, "The server's port configuration, default 8080")
	flag.Parse()

	if hp.port == "" {
		hp.port = "8080"
	}

	if hp.expression != "" {
		hp.code = CompileJQ(hp.expression)
		log.Debug().Msg("JQ expression loaded: " + hp.expression)
	}

	return hp
}

type HandlerParams struct {
	expression string
	code       *gojq.Code
	host       string
	port       string
	scheme     string
}
