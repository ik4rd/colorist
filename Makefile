BINARY  := colorist
PKG     := ./cmd/app
BIN_DIR := bin

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o $(BIN_DIR)/$(BINARY) $(PKG)

.PHONY: run
run:
	go run $(PKG) $(ARGS)

.PHONY: web
web:
	go run ./cmd/web $(ARGS)

# Desktop (Wails). Requires the wails CLI: see desktop/README.md.
# Builds a native app for the host OS into desktop/build/bin.
.PHONY: desktop
desktop:
	cd desktop && wails build $(ARGS)

# Live-reloading desktop dev window (rebuilds Go on change).
.PHONY: desktop-dev
desktop-dev:
	cd desktop && wails dev $(ARGS)

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out
