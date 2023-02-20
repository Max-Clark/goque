package main

import (
	"bytes"
	"flag"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const waitForServerRetries = 20
const waitForServerDelay = 50 * time.Millisecond

func _setEnvFromEnviron(environ []string) {
	os.Clearenv()
	for _, v := range environ {
		split := strings.Split(v, "=")
		os.Setenv(split[0], split[1])
	}
}

func _clearFlags() {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

func _setArgs(args []string) []string {
	oldArgs := os.Args
	os.Args = args
	return oldArgs
}

func _resetGetGoqueParamsFromStr(args []string) *GoqueParams {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = args

	_clearFlags()

	config := GetDefaultConfiguration()
	setEnvs, config := SetConfiguration(config)
	gp := ParseGoqueParams(setEnvs, config)
	return gp
}

// https://github.com/phayes/freeport/blob/master/freeport.go
func _getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func _waitForServer(client *http.Client, url string) error {
	var err error

	for i := 0; i < waitForServerRetries; i++ {
		time.Sleep(waitForServerDelay)
		_, err = client.Get(url)

		if err == nil {
			return nil
		}
	}

	return err
}

/****** TESTS ******/

func TestConfigDefaults(t *testing.T) {
	log.Info().Msg(reflect.Func.String())
	defaultConfig := GetDefaultConfiguration()
	_, config := SetConfiguration(GetDefaultConfiguration())

	for k, _ := range defaultConfig {
		assert.Equal(t, defaultConfig[k].val, config[k].val)
	}
}

func TestGetGoqueDefaults(t *testing.T) {
	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)

	assert.Equal(t, defaultEscapeHTML, gp.escape)
	assert.Equal(t, defaultHost, gp.host)
	assert.Equal(t, defaultPath, gp.path)
	assert.Equal(t, defaultPort, gp.port)
	assert.Equal(t, defaultScheme, gp.scheme)
	assert.Equal(t, defaultTracerDisable, gp.tracerDisabled)
	assert.Equal(t, defaultTracerRatio, gp.tracerRatio)
	assert.Equal(t, defaultTracerEndpoint, gp.tracerEndpoint)
}

// Set Env and command line args
func TestGetGoqueParamsEnvs(t *testing.T) {
	environ := os.Environ()
	defer _setEnvFromEnviron(environ)

	os.Clearenv()

	config := GetDefaultConfiguration()

	for k, v := range config {
		switch k {
		case "escapeHtml":
			os.Setenv(v.envVar, "true")
		case "tracerDisable":
			os.Setenv(v.envVar, "true")
		case "tracerRatio":
			os.Setenv(v.envVar, "0.3")
		default:
			os.Setenv(v.envVar, v.val+"1")
		}
	}

	args := []string{os.Args[0]}
	gp := _resetGetGoqueParamsFromStr(args)

	assert.Equal(t, !defaultEscapeHTML, gp.escape)
	assert.Equal(t, defaultHost+"1", gp.host)
	assert.Equal(t, defaultPath+"1", gp.path)
	assert.Equal(t, defaultPort+"1", gp.port)
	assert.Equal(t, defaultScheme+"1", gp.scheme)
	assert.Equal(t, !defaultTracerDisable, gp.tracerDisabled)
	assert.Equal(t, 0.3, gp.tracerRatio)
	assert.Equal(t, defaultTracerEndpoint+"1", gp.tracerEndpoint)
}

// Set Env and command line args, verify command line args overwrite
func TestGetGoqueParamsEnvsAndCli(t *testing.T) {
	environ := os.Environ()
	defer _setEnvFromEnviron(environ)

	os.Clearenv()

	config := GetDefaultConfiguration()

	for k, v := range config {
		switch k {
		case "escapeHtml":
			os.Setenv(v.envVar, "true")
		case "tracerDisable":
			os.Setenv(v.envVar, "true")
		case "tracerRatio":
			os.Setenv(v.envVar, "0.3")
		default:
			os.Setenv(v.envVar, v.val+"1")
		}
	}

	args := []string{os.Args[0]}

	for k, v := range config {
		args = append(args, "-"+v.arg)
		switch k {
		case "escapeHtml":
			args = append(args, "false")
		case "tracerDisable":
			args = append(args, "false")
		case "tracerRatio":
			args = append(args, "0.4")
		default:
			args = append(args, v.val+"2")
		}
	}

	gp := _resetGetGoqueParamsFromStr(args)

	assert.Equal(t, defaultEscapeHTML, gp.escape)
	assert.Equal(t, defaultHost+"2", gp.host)
	assert.Equal(t, defaultPath+"2", gp.path)
	assert.Equal(t, defaultPort+"2", gp.port)
	assert.Equal(t, defaultScheme+"2", gp.scheme)
	assert.Equal(t, defaultTracerDisable, gp.tracerDisabled)
	assert.Equal(t, 0.4, gp.tracerRatio)
	assert.Equal(t, defaultTracerEndpoint+"2", gp.tracerEndpoint)
}

func TestGetGoqueParamsBadParseToDefaults(t *testing.T) {
	environ := os.Environ()
	defer _setEnvFromEnviron(environ)

	os.Clearenv()

	config := GetDefaultConfiguration()

	for k, v := range config {
		switch k {
		case "escapeHtml":
			os.Setenv(v.envVar, "thisisnotabool")
		case "tracerDisable":
			os.Setenv(v.envVar, "thisisalsonotabool")
		case "tracerRatio":
			os.Setenv(v.envVar, "zeropointthree")
		case "logLevel":
			os.Setenv(v.envVar, "Plaid")
		default:
			os.Setenv(v.envVar, v.val)
		}
	}

	args := []string{os.Args[0]}

	gp := _resetGetGoqueParamsFromStr(args)

	assert.Equal(t, defaultEscapeHTML, gp.escape)
	assert.Equal(t, defaultHost, gp.host)
	assert.Equal(t, defaultPath, gp.path)
	assert.Equal(t, defaultPort, gp.port)
	assert.Equal(t, defaultScheme, gp.scheme)
	assert.Equal(t, defaultTracerDisable, gp.tracerDisabled)
	assert.Equal(t, defaultTracerRatio, gp.tracerRatio)
	assert.Equal(t, defaultTracerEndpoint, gp.tracerEndpoint)
}

func Test_main(t *testing.T) {
	tests := []struct {
		description    string
		route          string
		method         string
		reqContentType string
		reqBody        string
		reqJqHeader    string
		resCode        int
		resBody        string
	}{
		{
			description:    "Valid result",
			method:         "POST",
			route:          "http://0.0.0.0:8080" + defaultPath,
			reqContentType: "application/json",
			reqBody:        `{"test":{"peanuts": true,"pineapple":"nope."}}`,
			reqJqHeader:    `.test`,
			resCode:        200,
			resBody:        `{"peanuts": true,"pineapple":"nope."}`,
		},
		{
			description:    "Bad Body",
			method:         "POST",
			route:          "http://0.0.0.0:8080" + defaultPath,
			reqContentType: "application/json",
			reqBody:        `"test`,
			reqJqHeader:    `.`,
			resCode:        400,
			resBody:        `{"status":"error","message":"readStringSlowPath: unexpected end of input, error found in #5 byte of ...|\"test|..., bigger context ...|\"test|..."}`,
		},
	}

	json := jsoniter.Config{
		EscapeHTML: false,
	}.Froze()

	_clearFlags()
	_setArgs([]string{os.Args[0]})

	go main()

	client := &http.Client{}

	_waitForServer(client, "http://0.0.0.0:8080"+defaultPath)

	for _, test := range tests {
		var bodyReader io.Reader
		if test.reqBody != "" {
			bodyReader = bytes.NewReader([]byte(test.reqBody))
		}

		req, err := http.NewRequest(test.method, test.route, bodyReader)

		if err != nil {
			assert.FailNow(t, "Error creating http request")
		}

		req.Header.Add("x-goque-jq-filter", test.reqJqHeader)
		req.Header.Add("content-type", test.reqContentType)

		res, err := client.Do(req)

		if err != nil {
			assert.FailNow(t, "Error performing http request")
		}

		body, err := io.ReadAll(res.Body)
		bodyString := string(body)

		if err != nil {
			assert.FailNow(t, "Error reading request body: "+bodyString)
		}
		res.Body.Close()

		var bodyObj any
		json.Unmarshal(body, &bodyObj)

		var desiredObj any
		json.Unmarshal([]byte(test.resBody), &desiredObj)

		assert.Equalf(t, desiredObj, bodyObj, "Body did not match")

	}
}
