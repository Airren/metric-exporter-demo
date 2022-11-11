
.PHONY: build_mac
build_mac:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/exporter-demo main.go

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/exporter-demo main.go

.PHONY: build_image
build_image:
	docker build -f Dockerfile ./ -t metric-exporter-demo:latest

.PHONY: help
help:
	@echo 'build_mac        - build binary'
	@echo 'build_linux      - build linux binary'
	@echo 'build_image      - build docker images'