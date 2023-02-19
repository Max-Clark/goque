FROM golang:1.20 as build

WORKDIR /go/src/goque
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/goque ./cmd/goque

FROM gcr.io/distroless/base-debian11
COPY --from=build /go/bin/goque /
CMD ["/goque"]