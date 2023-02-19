/*
Goque is a high throughput HTTP JQ processor based on fiber and gojq.
Goque is highly configurable and made for both container and local usage.

Configuration of goque:

| Description           | Default        | Env Var     | CLI  | HTTP Header       |
| :-------------------- | :------------- | :---------- | :--- | :---------------- |
| JQ filter string      |                | JQ_FILTER   | -jq  | x-goque-jq-filter |
| JQ API path           | `"/api/v1/jq"` | JQ_PATH     | -a   |                   |
| Server host           | `""`           | HOST        | -h   |                   |
| Server port           | `"8080"`       | PORT        | -p   |                   |
| Escape HTML on return | `false`        | HTML_ESCAPE | -e   |                   |

Usage of ./goque:
  -a string
        Server path (default "/api/v1/jq")
  -e string
        Escape HTML on return (default "false")
  -h string
        Server host
  -jq string
        JQ filter string
  -l string
        Default log level (default "Info")
  -p string
        Server port (default "8080")
  -s string
        Server scheme
  -td string
        Disable tracer (default "false")
  -te string
        Tracer endpoint, url (default "http://localhost:14268/api/traces")
  -tr string
        Tracer ratio, 0-1 (default "1")
*/

package main

import (
	"context"
	"flag"
	"os"
	"strconv"

	"github.com/itchyny/gojq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultLogLevel = zerolog.DebugLevel
const hostDefault = ""
const portDefault = "8080"
const pathDefault = "/api/v1/jq"
const escapeHTMLDefault = "false"
const schemeDefault = ""
const tracerEndpointDefault = "http://localhost:14268/api/traces"

// Entry to goque. Initializes logger, gets params, and
// starts server
func main() {
	gp := GetGoqueParams()

	InitLogging(gp.logLevel)

	tp := InitTracer(gp.tracerRatio, gp.tracerEndpoint)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	RunServer(gp)
}

// Grabs the handler params from the env vars or command line.
// Command line has precedence. Returns a HandlerParams object
// that configures the server & jq evaluation.
func GetGoqueParams() *GoqueParams {

	var config = map[string]*ConfigurationVar{
		"jq":             {desc: "JQ filter string", val: "", envVar: "GOQUE_JQ_FILTER", arg: "jq"},
		"path":           {desc: "Server path", val: pathDefault, envVar: "GOQUE_PATH", arg: "a"},
		"host":           {desc: "Server host", val: hostDefault, envVar: "GOQUE_HOST", arg: "h"},
		"port":           {desc: "Server port", val: portDefault, envVar: "GOQUE_PORT", arg: "p"},
		"scheme":         {desc: "Server scheme", val: schemeDefault, envVar: "GOQUE_SCHEME", arg: "s"},
		"escapeHtml":     {desc: "Escape HTML on return", val: escapeHTMLDefault, envVar: "GOQUE_HTML_ESCAPE", arg: "e"},
		"logLevel":       {desc: "Default log level", val: "Info", envVar: "GOQUE_LOG_LEVEL", arg: "l"},
		"tracerDisable":  {desc: "Disable tracer", val: "false", envVar: "GOQUE_TRACER_DISABLE", arg: "td"},
		"tracerRatio":    {desc: "Tracer ratio, 0-1", val: "1", envVar: "GOQUE_TRACER_RATIO", arg: "tr"},
		"tracerEndpoint": {desc: "Tracer endpoint, url", val: tracerEndpointDefault, envVar: "GOQUE_TRACER_ENDPOINT", arg: "te"},
	}

	for _, v := range config {
		// If the env var exists, set the value
		if v.envVar != "" {
			if envVal, ok := os.LookupEnv(v.envVar); ok {
				v.val = envVal
			}
		}

		// If the command line argument was set, set the value.
		// Note that this overwrites the env vars.
		flag.StringVar(&v.val, v.arg, v.val, v.desc)
	}

	// Load the flags
	flag.Parse()

	// escape is a boolean, so parse it. If error, set to false and warn.
	parsedEscapeHtml, err := strconv.ParseBool(config["escapeHtml"].val)
	if err != nil {
		log.Warn().Msg("-e or HTML_ESCAPE invalid, defaulting to `false`")
		parsedEscapeHtml = false
	}

	tracerRatio, err := strconv.ParseFloat(config["tracerRatio"].val, 64)
	if err != nil {
		log.Warn().AnErr("InitTracer", err).Msg("Invalid tracerRatio, defaulting to `1` (100%)")
		tracerRatio = 1.0
	}

	logLevel, err := zerolog.ParseLevel(config["logLevel"].val)
	if err != nil {
		log.Warn().AnErr("InitTracer", err).Msg("Invalid logLevel, defaulting to " + defaultLogLevel.String())
		logLevel = defaultLogLevel
	}

	return &GoqueParams{
		tracerEnabled:  true,
		tracerRatio:    tracerRatio,
		tracerEndpoint: config["tracerEndpoint"].val,
		logLevel:       logLevel,
		jqFilter:       config["jq"].val,
		host:           config["host"].val,
		port:           config["port"].val,
		path:           config["path"].val,
		scheme:         config["scheme"].val,
		escape:         parsedEscapeHtml,
	}
}

// TODO
func PrintGoqueParams(p GoqueParams) {

}

type ConfigurationVar struct {
	arg    string
	desc   string
	envVar string
	val    string
}

// A struct containing server and jq configuration info.
type GoqueParams struct {
	code           *gojq.Code // Compiled JQ if set with env/cli
	tracerEnabled  bool
	tracerRatio    float64
	tracerEndpoint string
	logLevel       zerolog.Level
	escape         bool   // Escape HTML
	jqFilter       string // The JQ filter string
	host           string // The server host
	port           string // The server port
	scheme         string // The server scheme
	path           string // The jq API path
}
