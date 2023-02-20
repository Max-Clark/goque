package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracerSmokeTest(t *testing.T) {
	tracerProvider := InitTracer(1, "http://localhost:14268/api/traces")
	assert.NotNil(t, tracerProvider)

	tracer := tracerProvider.Tracer("")

	_, span := tracer.Start(context.Background(), "test")
	defer span.End()
}
