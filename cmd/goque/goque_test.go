package main

import (
	"bytes"
	"flag"
	"io"
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

// Waits for the a server to be ready (i.e., return a HTTP response)
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

	// Configure the JSON parser
	json := jsoniter.Config{
		EscapeHTML: false,
	}.Froze()

	os.Clearenv()
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

func TestConfigDefaults(t *testing.T) {
	_clearFlags()
	os.Clearenv()

	log.Info().Msg(reflect.Func.String())
	defaultConfig := GetDefaultConfiguration()
	_, config := SetConfiguration(GetDefaultConfiguration())

	for k, _ := range defaultConfig {
		assert.Equal(t, defaultConfig[k].val, config[k].val)
	}
}

func TestParseGoqueParams(t *testing.T) {
	type args struct {
		setEnvs map[string]string
		config  map[string]*ConfigurationVar
	}

	test2Prep := GetDefaultConfiguration()
	for k, v := range test2Prep {
		if k != "jq" {
			v.val = v.envVar
		}
	}

	test3Prep := GetDefaultConfiguration()
	for k, v := range test3Prep {
		if k != "jq" {
			v.val = v.envVar
		} else {
			v.val = "."
		}
	}

	tests := []struct {
		name string
		args args
		want *GoqueParams
	}{
		{
			name: "test1 - defaults",
			args: args{
				setEnvs: map[string]string{},
				config:  GetDefaultConfiguration(),
			},
			want: &GoqueParams{
				code:           nil,
				tracerDisabled: defaultTracerDisable,
				tracerRatio:    defaultTracerRatio,
				tracerEndpoint: defaultTracerEndpoint,
				escape:         defaultEscapeHTML,
				host:           defaultHost,
				port:           defaultPort,
				scheme:         defaultScheme,
				path:           defaultPath,
			},
		},
		{
			name: "test2 - parse errors (use defaults)",
			args: args{
				setEnvs: map[string]string{},
				config:  test2Prep,
			},
			want: &GoqueParams{
				code:           nil,
				tracerDisabled: false,
				tracerRatio:    1.0,
				tracerEndpoint: "GOQUE_TRACER_ENDPOINT",
				escape:         false,
				host:           "GOQUE_HOST",
				port:           "GOQUE_PORT",
				scheme:         "GOQUE_SCHEME",
				path:           "GOQUE_PATH",
			},
		},
		{
			name: "test3 - compile jq",
			args: args{
				setEnvs: map[string]string{},
				config:  test3Prep,
			},
			want: &GoqueParams{
				code:           CompileJQCode("."),
				tracerDisabled: false,
				tracerRatio:    1.0,
				tracerEndpoint: "GOQUE_TRACER_ENDPOINT",
				escape:         false,
				host:           "GOQUE_HOST",
				port:           "GOQUE_PORT",
				scheme:         "GOQUE_SCHEME",
				path:           "GOQUE_PATH",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseGoqueParams(tt.args.setEnvs, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGoqueParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetConfiguration(t *testing.T) {
	type args struct {
		config  map[string]*ConfigurationVar
		envVars map[string]string
		flags   []string
	}

	test2Configuration := GetDefaultConfiguration()
	test2EnvVars := make(map[string]string)
	for k, v := range test2Configuration {
		test2EnvVars[v.envVar] = v.val + "2"
		test2Configuration[k].val += "2"
	}

	tests := []struct {
		name  string
		args  args
		want  map[string]string
		want1 map[string]*ConfigurationVar
	}{
		{
			name: "test 1 - defaults",
			args: args{
				config:  GetDefaultConfiguration(),
				envVars: make(map[string]string),
				flags:   []string{os.Args[0]},
			},
			want1: GetDefaultConfiguration(),
			want:  make(map[string]string),
		},
		{
			name: "test 2 - envs",
			args: args{
				config:  GetDefaultConfiguration(),
				envVars: test2EnvVars,
				flags:   []string{os.Args[0]},
			},
			want1: test2Configuration,
			want:  test2EnvVars,
		},
	}
	for _, tt := range tests {
		environ := os.Environ()
		defer _setEnvFromEnviron(environ)

		_clearFlags()
		os.Clearenv()

		for k, v := range tt.args.envVars {
			os.Setenv(k, v)
		}

		_setArgs(tt.args.flags)

		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SetConfiguration(tt.args.config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetConfiguration() got = %v, want %v", got, tt.want)
			}

			assert.Len(t, got1, len(tt.args.config))

			for k, _ := range got1 {
				if !reflect.DeepEqual(got1[k], tt.want1[k]) {
					t.Errorf("SetConfiguration() got1 = %v, want %v", got1[k], tt.want1[k])
				}
			}
		})
	}
}
