#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o build/paperless-ngx-copier_linux_x86 paperless-ngx-copier.go
GOOS=linux GOARCH=arm64 go build -o build/paperless-ngx-copier_linux_arm64 paperless-ngx-copier.go

GOOS=windows GOARCH=amd64 go build -o build/paperless-ngx-copier_win_x86.exe paperless-ngx-copier.go
GOOS=windows GOARCH=arm64 go build -o build/paperless-ngx-copier_win_arm64.exe paperless-ngx-copier.go

GOOS=darwin GOARCH=amd64 go build -o build/paperless-ngx-copier_darwin_x86 paperless-ngx-copier.go
GOOS=darwin GOARCH=arm64 go build -o build/paperless-ngx-copier_darwin_arm64 paperless-ngx-copier.go