package main

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

// Starts the http server. Takes params for escaping html,
// server properties, and other handler variables. Handles
// json POSTs on hp.path.
func RunServer(hp *HandlerParams) {
	json := jsoniter.Config{
		EscapeHTML: hp.escape,
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

	// TODO: Host OAS spec

	app.Post(hp.path, func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		c.AcceptsCharsets("utf-8")
		return PostHandler(c, hp)
	})

	var parsedPort = hp.port
	if parsedPort != "" {
		parsedPort = ":" + parsedPort
	}

	var parsedScheme = hp.scheme
	if parsedScheme != "" {
		parsedScheme = hp.scheme + "//"
	}

	var url = parsedScheme + hp.host + parsedPort

	log.Fatal().AnErr("RunServer", app.Listen(url)).Msg("")
}
