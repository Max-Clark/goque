package main

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

// Starts the http server. Takes params for escaping html,
// server properties, and other handler variables. Handles
// json POSTs on hp.path.
func RunServer(gp *GoqueParams) {
	json := jsoniter.Config{
		EscapeHTML: gp.escape,
	}.Froze()

	app := fiber.New(fiber.Config{
		AppName:               "goque",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	// }))

	app.Use(otelfiber.Middleware())

	// TODO: Host OAS spec

	app.Post(gp.path, func(c *fiber.Ctx) error {
		_, span := tracer.Start(c.UserContext(), "PostHandler")
		defer span.End()
		return HandlePost(c, gp)
	})

	// TODO: Create better validation for server urls
	// :<port> is supported
	// [<scheme>//]<host>[:<port>] is supported
	var parsedPort = gp.port
	if parsedPort != "" {
		parsedPort = ":" + parsedPort
	}

	var parsedScheme = gp.scheme
	if parsedScheme != "" {
		parsedScheme = gp.scheme + "//"
	}

	var url = parsedScheme + gp.host + parsedPort

	log.Fatal().AnErr("RunServer", app.Listen(url)).Msg("")
}
