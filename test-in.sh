#!/usr/bin/env sh

go build ./cmd/in

echo \
'{
  "source": {
    "defectdojo_url": "https://something",
    "username": "admin",
    "api_key": "key",
    "debug": true
  },
  "params": {}
}' | ./in

rm in
