FROM golang:1.20 as build

WORKDIR /go/src/goque
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/goque ./cmd/goque

FROM gcr.io/distroless/base-debian11

LABEL org.opencontainers.image.source=https://github.com/Max-Clark/goque
LABEL org.opencontainers.image.description="A blazing fast and dead simple http-based jq processor written in go"
LABEL org.opencontainers.image.licenses=MIT

COPY --from=build /go/bin/goque /
CMD ["/goque"]