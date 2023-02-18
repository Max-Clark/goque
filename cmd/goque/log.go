package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogging(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(level)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Info().Msg(name + " took " + elapsed.String())
}
