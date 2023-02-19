# goque

A blazing fast and dead simple http-based jq processor written in go.


## Features

- Fast HTTP service with [Fiber](https://gofiber.io/)
- Fast JQ processing with [gojq](https://github.com/itchyny/gojq)
- Jaeger metrics
- A playground for testing TBD
- Small footprint

### Tracing 

Goque exports jaeger metrics. See tracing configuration.

## Installation

### Local Copy
- `git clone https://github.com/Max-Clark/goque.git`
- `cd ./goque`
- `go build ./cmd/goque`

### Go Installation

TBD

### Docker

- `docker build --tag local/goque -f goque.dockerfile .`
- `docker run -d --name=goque -p 8888:8080 local/goque`

## Usage

Goque is highly configurable, but defaults will work for most deployments.

```sh
./goque  # Start a new goque instance
```

```sh
curl --request POST \
  --url http://localhost:8080/api/v1/jq \
  --header 'Content-Type: application/json' \
  --header 'x-goque-jq-filter: .test' \
  --data '{
        "test": {
                "peanuts": true,
                "pineapple": "nope."
        }
}'
{"peanuts":true,"pineapple":"nope."}%
```

Assigning a JQ filter with environment variables or command line will compile
the jq code, resulting in faster processing.

```sh
# GOQUE_JQ_FILTER='."test"' ./goque
./goque -jq '."test"'  # Both work, but cli has preference
```

### Goque Configuration

*NOTE* Variable preference is Env Var < Command Line < HTTP Header

| Description           | Default                             | Env Var                  | CLI  | HTTP Header       |
| :-------------------- | :---------------------------------- | :----------------------- | :--- | :---------------- |
| JQ filter string      | `nil`                               | GOQUE_JQ_FILTER          | -jq  | x-goque-jq-filter |
| JQ API path           | `"/api/v1/jq"`                      | GOQUE_PATH               | -a   |                   |
| Server host           | `""`                                | GOQUE_HOST               | -h   |                   |
| Server port           | `"8080"`                            | GOQUE_PORT               | -p   |                   |
| Escape HTML on return | `false`                             | GOQUE_HTML_ESCAPE        | -e   |                   |
| Default log level     | `Info`                              | GOQUE_LOG_LEVEL          | -l   |                   |
| Tracer disable        | `false`                             | GOQUE_TRACER_DISABLE     | -td  |                   |
| Tracer ratio, \[0,1\] | `1`                                 | GOQUE_TRACER_RATIO       | -tr  |                   |
| Tracer export dest.   | `http://localhost:14268/api/traces` | GOQUE_TRACER_EXPORT_DEST | -te  |                   |

## Building 

go build -v -ldflags="-X 'main.Version=v1.0.0' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'"

---

# Plan

- [x] Initial commit
- [x] Logging
    - [x] Research logging
    - [x] Implement basic logging
    - [x] More robust logging
- [x] Tracing
- [ ] Metrics
- ~[ ] OAS~ It's a pretty obvious API, going to work on gojqplay instead
- [-] Configuration Validation
  - [x] Crucial items validated
- JQ
  - [x] Investigate gojq
      - Yup, it's fast
  - [x] Implement gojq from env variable
      - [x] Compile filter
      - [x] Error on bad filter
  - [x] Implement gojq from http header
      - [x] Error on bad filterr
  - [ ] Implement benchmarking scaffolding 
  - [ ] Research testing methodologies/libraries
  - [ ] Implement testing scaffolding
- HTTP
  - [x] Implement basic server with `http`
  - [x] Investigate http libraries
  - [ ] Implement TLS
  - [ ] Investigate websocket usage
  - [ ] Investigate sidecar usage
      - [ ] Proper implementation? MITM?
  - [ ] Implement benchmarking scaffolding 
  - [ ] Research testing methodologies/libraries
  - [ ] Implement testing scaffolding
- ~GoquePlay~ Going to make gojqplay instead, probably better to separate and probably more popular
- CI
  - [ ] Lock main branch, merge by request
  - [ ] Run tests
  - [ ] Run vuln scan
  - [ ] Build wasm module
  - [ ] Test wasm module
  - [ ] Github release
  - [ ] Build image and push
    