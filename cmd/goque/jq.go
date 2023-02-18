package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itchyny/gojq"
	"github.com/rs/zerolog/log"
)

func CompileJQ(expression string) *gojq.Code {
	query, err := gojq.Parse(expression)
	if err != nil {
		log.Fatal().AnErr("JQ", err).Msg("An invalid JQ expression was entered")
	}

	code, err := gojq.Compile(query)
	if err != nil {
		log.Fatal().AnErr("JQ", err).Msg("Could not compile JQ expression")
	}

	return code
}

func RunCompiled(input any, code *gojq.Code) (interface{}, error) {
	iter := code.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		if v != nil {
			return v, nil
		}
	}
	return nil, nil
}

func PostHandler(c *fiber.Ctx, p *HandlerParams) error {
	// defer timeTrack(time.Now(), "PostHandler")

	if p != nil && p.code != nil {
		var body interface{}
		err := c.BodyParser(&body)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		out, err := RunCompiled(body, p.code)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(out)
	}

	c.Status(fiber.StatusBadRequest)
	return c.JSON(fiber.Map{"status": "error", "message": "A JQ expression was not sent with request"})
}
