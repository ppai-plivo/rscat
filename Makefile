.PHONY: build release

build:
	CGO_ENABLED=0 go build -o rscat

release:
	CGO_ENABLED=0 go build -mod=readonly -ldflags "-w -s" -o rscat
