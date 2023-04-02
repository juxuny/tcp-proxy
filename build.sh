#!/bin/bash
GOOS=linux GOHOSTARCH=amd64 GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/tcp-proxy_linux_x64_v1.1.3