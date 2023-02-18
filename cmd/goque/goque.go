package main

import (
	"flag"
	"os"

	"github.com/itchyny/gojq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultLogLevel = zerolog.DebugLevel
const hostDefault = ""
const portDefault = "8080"
const pathDefault = "/api/v1/jq"

func main() {
	InitLogging(defaultLogLevel)

	hp := GetHandlerParams()

	RunServer(&hp)
}

func GetHandlerParams() HandlerParams {
	hp := HandlerParams{
		jqFilter: os.Getenv("JQ_FILTER"),
		host:     os.Getenv("HOST"),
		port:     os.Getenv("PORT"),
		path:     os.Getenv("JQ_PATH"),
		scheme:   "", // TODO: https
	}

	flag.StringVar(&hp.jqFilter, "jq", hp.jqFilter, "The jq filter")
	flag.StringVar(&hp.host, "h", hp.host, "The server's host configuration, default "+hostDefault)
	flag.StringVar(&hp.port, "p", hp.port, "The server's port configuration, default "+portDefault)
	flag.StringVar(&hp.port, "a", hp.port, "The server's path configuration, default "+pathDefault)
	flag.Parse()

	if hp.port == "" {
		hp.port = "8080"
	}

	if hp.jqFilter != "" {
		hp.code = CompileJQ(hp.jqFilter)
		log.Debug().Msg("JQ filter loaded: " + hp.jqFilter)
	}

	return hp
}

type HandlerParams struct {
	code     *gojq.Code
	jqFilter string
	host     string
	port     string
	scheme   string
	path     string
}
