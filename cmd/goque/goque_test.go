package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _setEnvFromEnviron(environ []string) {
	os.Clearenv()
	for _, v := range environ {
		split := strings.Split(v, "=")
		os.Setenv(split[0], split[1])
	}
}

func _resetGetGoqueParamsFromStr(args []string) *GoqueParams {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = args

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	config := GetDefaultConfiguration()
	setEnvs, config := SetConfiguration(config)
	gp := ParseGoqueParams(setEnvs, config)
	return gp
}

func TestConfigDefaults(t *testing.T) {
	defaultConfig := GetDefaultConfiguration()
	_, config := SetConfiguration(nil)

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
