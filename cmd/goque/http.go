package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RunServer(hp *HandlerParams) {
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

	app.Post("/api/v1/jq", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
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

	log.Fatal().Err(app.Listen(url))
}
