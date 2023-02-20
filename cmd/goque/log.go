package main

import (
	"github.com/rs/zerolog"
)

// Initialize logging. Setting level will set the default logging level.
func InitLogging(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(level)
}

// A tool for tracing function time. To use, defer timeTrack at the
// top of the function to track.
// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Info().Msg(name + " took " + elapsed.String())
// }
