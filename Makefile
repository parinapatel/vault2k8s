LDFLAGS := -s -w

.PHONY: all
all: build

bin:
	mkdir -p bin/

.PHONY: build
build: bin
	mkdir -p bin/linux-amd64/
	mkdir -p bin/OSX/
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/linux-amd64/vault2k8s -v main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/OSX/vault2k8s -v main.go

.PHONY: clean
clean:
	rm -rf bin/
	go clean
