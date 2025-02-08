.DEFAULT_GOAL := build
.PHONY:fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

clean:
	del gopherhole.exe

run: clean build
	gopherhole