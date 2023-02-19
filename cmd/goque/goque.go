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

Usage of goque:
  -a string
        Server path configuration; default "/api/v1/jq"
  -e    Escape HTML on return, use when returning to a web interface; default false
  -h string
        Server host configuration; default ""
  -jq string
        JQ filter string
  -p string
        Server port configuration; default "8080"
*/

package main

import (
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

// Entry to goque. Initializes logger, gets params, and
// starts server
func main() {
	InitLogging(defaultLogLevel)

	hp := GetHandlerParams()

	RunServer(&hp)
}

// Grabs the handler params from the env vars or command line.
// Command line has precedence. Returns a HandlerParams object
// that configures the server & jq evaluation.
func GetHandlerParams() HandlerParams {

	var config = map[string]*ConfigurationVar{
		"jq":         {desc: "JQ filter string", val: "", envVar: "JQ_FILTER", arg: "jq"},
		"path":       {desc: "Server path", val: pathDefault, envVar: "JQ_PATH", arg: "a"},
		"host":       {desc: "Server host", val: hostDefault, envVar: "HOST", arg: "h"},
		"port":       {desc: "Server port", val: portDefault, envVar: "PORT", arg: "p"},
		"scheme":     {desc: "Server scheme", val: schemeDefault, envVar: "SCHEME", arg: "s"},
		"escapeHtml": {desc: "Escape HTML on return", val: escapeHTMLDefault, envVar: "HTML_ESCAPE", arg: "e"},
	}

	for _, v := range config {
		// If the env var exists, set the value
		if envVal, ok := os.LookupEnv(v.envVar); ok {
			v.val = envVal
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
		log.Warn().Msg("-e or HTML_ESCAPE invalid, defaulting to false")
	}

	return HandlerParams{
		jqFilter: config["jq"].val,
		host:     config["host"].val,
		port:     config["port"].val,
		path:     config["path"].val,
		scheme:   config["scheme"].val,
		escape:   parsedEscapeHtml,
	}
}

type ConfigurationVar struct {
	arg    string
	desc   string
	envVar string
	val    string
}

// A struct containing server and jq configuration info.
type HandlerParams struct {
	code     *gojq.Code // Compiled JQ if set with env/cli
	escape   bool       // Escape HTML
	jqFilter string     // The JQ filter string
	host     string     // The server host
	port     string     // The server port
	scheme   string     // The server scheme
	path     string     // The jq API path
}
