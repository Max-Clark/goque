# goque

A blazing fast http jq evaluator written in go.


## HTTP Server

Goque uses [Fiber](https://gofiber.io/) as its server/router. 

## Usage

### Goque Configuration

*NOTE* Variable preference is Env Var < Command Line < HTTP Header

| Description           | Default        | Env Var           | CLI  | HTTP Header       |
| :-------------------- | :------------- | :---------------- | :--- | :---------------- |
| JQ filter string      | `nil`          | GOQUE_JQ_FILTER   | -jq  | x-goque-jq-filter |
| JQ API path           | `"/api/v1/jq"` | GOQUE_PATH        | -a   |                   |
| Server host           | `""`           | GOQUE_HOST        | -h   |                   |
| Server port           | `"8080"`       | GOQUE_PORT        | -p   |                   |
| Escape HTML on return | `false`        | GOQUE_HTML_ESCAPE | -e   |                   |
| Default log level     | `Info`         | GOQUE_LOG_LEVEL   | -l   |                   |


### 3rd Party Configuration

#### Opencollector Tracing 

| Description               | Default | Env Var              | CLI  |
| :------------------------ | :------ | :------------------- | :--- |
| Disable tracer            | `false` | GOQUE_TRACER_DISABLE | -td  |
| Set tracer ratio, \[0,1\] | `1`     | GOQUE_TRACER_RATIO   | -tr  |



---

# Plan

- [x] Initial commit
- [ ] Logging/metrics/tracing
    - [x] Research logging
    - [x] Implement basic logging
    - [x] More robust logging
- [ ] Tracing
- [ ] Metrics
- [ ] OAS
- [ ] Validation
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
    