package main

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type JQMatch struct {
	filter string
	input  string
	output string
}

func _GetNewFiberContext() *fiber.Ctx {
	app := fiber.New()
	return app.AcquireCtx(&fasthttp.RequestCtx{})
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

func TestHandlerNoConfiguredJQ(t *testing.T) {
	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)
	c := _GetNewFiberContext()

	c.Context().Request.SetBody([]byte(`{"peanuts":true,"pineapple":"nope."}`))
	c.Context().Request.Header.Add("content-type", "application/json")

	HandlePost(c, gp)

	body := string(c.Response().Body())

	assert.Equal(t, "{\"message\":\"A JQ filter was not sent with request\",\"status\":\"error\"}", body)
	assert.Equal(t, c.Response().StatusCode(), fiber.StatusBadRequest)
}

func TestHandlerImproperJQFilter(t *testing.T) {
	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)
	c := _GetNewFiberContext()

	c.Context().Request.SetBody([]byte(`{"peanuts":true,"pineapple":"nope."}`))
	c.Context().Request.Header.Add("content-type", "application/json")
	c.Context().Request.Header.Add("x-goque-jq-filter", "(wut")

	assert.NoError(t, HandlePost(c, gp))

	body := string(c.Response().Body())

	assert.Equal(t, "{\"message\":\"unexpected EOF\",\"status\":\"error\"}", body)
	assert.Equal(t, c.Response().StatusCode(), fiber.StatusBadRequest)
}

func TestHandlerSuccessHeader(t *testing.T) {
	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)
	c := _GetNewFiberContext()

	c.Context().Request.SetBody([]byte(`{"peanuts":true,"pineapple":"nope."}`))
	c.Context().Request.Header.Add("content-type", "application/json")
	c.Context().Request.Header.Add("x-goque-jq-filter", ".peanuts")

	assert.NoError(t, HandlePost(c, gp))

	body := string(c.Response().Body())

	assert.Equal(t, "true", body)
	assert.Equal(t, c.Response().StatusCode(), fiber.StatusOK)
}

func TestHandlerSuccessCompiled(t *testing.T) {
	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)
	gp.code = CompileJQCode(".peanuts")
	c := _GetNewFiberContext()

	c.Context().Request.SetBody([]byte(`{"peanuts":true,"pineapple":"nope."}`))
	c.Context().Request.Header.Add("content-type", "application/json")

	assert.NoError(t, HandlePost(c, gp))

	body := string(c.Response().Body())

	assert.Equal(t, "true", body)
	assert.Equal(t, c.Response().StatusCode(), fiber.StatusOK)
}
