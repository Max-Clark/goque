# goque

A fast http jq evaluator

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
    - [ ] Implement gojq from http header
        - [ ] Error on bad filterr
    - [ ] Implement benchmarking scaffolding 
    - [ ] Research testing methodologies/libraries
    - [ ] Implement testing scaffolding
- HTTP
    - [x] Implement basic server with `http`
    - [x] Investigate http libraries
    - [ ] Investigate websocket usage
    - [ ] Investigate sidecar usage
        - [ ] Proper implementation? MITM?
    - [ ] Implement benchmarking scaffolding 
    - [ ] Research testing methodologies/libraries
    - [ ] Implement testing scaffolding
    