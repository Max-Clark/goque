package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itchyny/gojq"
	"github.com/rs/zerolog/log"
)

// Compile the provided filter from env vars. Failing the parse or
// compile will fatal the program.
func CompileJQCode(filter string) *gojq.Code {
	query, err := gojq.Parse(filter)
	if err != nil {
		log.Fatal().AnErr("JQ", err).Msg("An invalid JQ filter was entered")
	}

	code, err := gojq.Compile(query)
	if err != nil {
		log.Fatal().AnErr("JQ", err).Msg("Could not compile JQ filter")
	}

	return code
}

// Returns the first value in the iter.
func GetFirstValueIter(iter gojq.Iter) (any, error) {
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

// The handler for jq evaluation requests. If a jq filter was provided
// with env vars, the resultant compiled code will be used here. If
// the x-goque-jq-filter header is set, the filter is parsed and ran
// against the body. The x-goque-jq-filter takes priority over
// JQ_FILTER. Parsing errors will be returned as 400 with a reason.
func PostHandler(c *fiber.Ctx, p *HandlerParams) error {
	// defer timeTrack(time.Now(), "PostHandler")

	if p == nil {
		log.Panic().Msg("HandlerParams not configured, panic")
	}

	// Parse the JSON body into object
	var body interface{}
	err := c.BodyParser(&body)

	// 400 if bad body
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// If jq filter header is set, prioritize over compiled code
	if jqHeader := c.Get("x-goque-jq-filter"); jqHeader != "" {
		query, err := gojq.Parse(jqHeader)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		out, err := GetFirstValueIter(query.Run(body))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(out)
	}

	// If env jq query was compiled run the query
	if p.code != nil {
		out, err := GetFirstValueIter(p.code.Run(body))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.JSON(out)
	}

	// jq filter nor jq env variable was provided
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "A JQ filter was not sent with request"})
}
