#!/usr/bin/env sh

go build ./cmd/out

echo \
'{
  "source": {
    "defectdojo_url": "https://something",
    "username": "admin",
    "api_key": "key",
    "debug": true
  },
  "params": {
    "report_type": "ZAP Scan"
  }
}' | ./out

rm out
