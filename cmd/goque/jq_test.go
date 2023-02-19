package main

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

type JQMatch struct {
	filter string
	input  string
	output string
}

func _CompileAssertEqual(t *testing.T, jm JQMatch) {
	json := jsoniter.Config{
		EscapeHTML: false,
	}.Froze()

	code := CompileJQCode(jm.filter)
	assert.NotNil(t, code)

	var inputBytes = []byte(jm.input)
	var data interface{}
	json.Unmarshal(inputBytes, &data)

	iter := code.Run(data)
	assert.NotNil(t, iter)

	out, ok := iter.Next()
	assert.True(t, ok)

	if err, ok := out.(error); ok {
		assert.NoError(t, err)
	}

	outBytes, err := json.Marshal(out)

	assert.NoError(t, err)

	outString := string(outBytes)
	assert.Equal(t, jm.output, outString)
}

func TestCompileJQ(t *testing.T) {

	// TODO: Test fatal?

	/* Valid Entries */
	_CompileAssertEqual(t, JQMatch{filter: ".", input: "\"test\"", output: "\"test\""})
	_CompileAssertEqual(t, JQMatch{filter: ".test", input: "{\"test\":\"2\"}", output: "\"2\""})
	_CompileAssertEqual(t, JQMatch{filter: ".test", input: "{\"test\":{\"value\":\"2\"}}", output: "{\"value\":\"2\"}"})
	_CompileAssertEqual(t, JQMatch{filter: ".test", input: "{\"test\":{\"value\":\"&\"}}", output: "{\"value\":\"&\"}"})
	_CompileAssertEqual(t, JQMatch{filter: ".test", input: "{\"test\":{\"value\":\"😍\"}}", output: "{\"value\":\"😍\"}"})
	_CompileAssertEqual(t, JQMatch{filter: ".\"日本\"", input: "{\"日本\":{\"value\":\"円\"}}", output: "{\"value\":\"円\"}"})

}
