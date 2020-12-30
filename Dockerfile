FROM golang:alpine as builder
WORKDIR /resource
COPY . /resource

ENV CGO_ENABLED 0
RUN go test ./client ./pkg  -failfast
RUN mkdir /resource/bin
# RUN go build -o /resource/bin/check ./cmd/check
RUN go build -o /resource/bin/in ./cmd/in
RUN go build -o /resource/bin/out ./cmd/out

FROM alpine:edge AS resource
RUN apk add --no-cach bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder /resource/bin /opt/resource
RUN chmod +x /opt/resource/*

# Test binaries exist
RUN stat /opt/resource/out