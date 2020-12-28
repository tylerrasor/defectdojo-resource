FROM golang:alpine as builder
WORKDIR /resource
COPY . /resource

ENV CGO_ENABLED 0
RUN go test -failfast
# RUN go build -o /resource/check ./cmd/check
# RUN go build -o /resource/in ./cmd/in
RUN go build -o /resource/out ./cmd/out

FROM alpine:edge AS resource
RUN apk add --no-cach bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder /resource /opt/resource
RUN chmod +x /opt/resource/*

# Test binaries exist
RUN stat /opt/resource/out