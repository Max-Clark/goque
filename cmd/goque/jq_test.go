package main

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

type JQMatch struct {
	filter string
	input  string
	output string
}

func ValidJqMatches() []JQMatch {
	return []JQMatch{
		{filter: ".", input: "\"test\"", output: "\"test\""},
		{filter: ".test", input: "{\"test\":\"2\"}", output: "\"2\""},
		{filter: ".test", input: "{\"test\":{\"value\":\"2\"}}", output: "{\"value\":\"2\"}"},
	}
}

func TestCompileJQ(t *testing.T) {
	// TODO: Test fatal?
	for _, v := range ValidJqMatches() {
		code := CompileJQ(v.filter)
		assert.NotNil(t, code)

		var jsonInput any
		json.Unmarshal([]byte(v.input), &jsonInput)

		iter := code.Run(jsonInput)
		assert.NotNil(t, iter)

		out, ok := iter.Next()
		assert.True(t, ok)

		if err, ok := out.(error); ok {
			assert.NoError(t, err)
		}

		outString, err := json.Marshal(out)

		// err, ok = out.(error)
		// assert.True(t, ok)
		assert.NoError(t, err)

		assert.Equal(t, v.output, string(outString))
	}
}
