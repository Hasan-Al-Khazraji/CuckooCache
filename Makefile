# Go Makefile for CuckooCache
# Builds binaries into bin/ and provides common utility targets

GO := go
BIN_DIR := bin

.PHONY: all
all: build

.PHONY: build
build: build-orchestrator build-worker build-ccctl

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: build-orchestrator
build-orchestrator: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/orchestrator ./cmd/orchestrator

.PHONY: build-worker
build-worker: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/worker ./cmd/worker

.PHONY: build-ccctl
build-ccctl: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/ccctl ./cmd/ccctl

.PHONY: run-orchestrator
run-orchestrator: build-orchestrator
	$(BIN_DIR)/orchestrator

.PHONY: run-worker
run-worker: build-worker
	$(BIN_DIR)/worker

.PHONY: run-ccctl
run-ccctl: build-ccctl
	$(BIN_DIR)/ccctl

.PHONY: test
test:
	$(GO) test ./...

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: help
help:
	@echo "Targets:"
	@echo "  all                - Build all binaries (default)"
	@echo "  build              - Build orchestrator, worker, ccctl"
	@echo "  build-orchestrator - Build orchestrator binary"
	@echo "  build-worker       - Build worker binary"
	@echo "  build-ccctl        - Build ccctl binary"
	@echo "  run-orchestrator   - Build and run orchestrator"
	@echo "  run-worker         - Build and run worker"
	@echo "  run-ccctl          - Build and run ccctl"
	@echo "  test               - Run go tests"
	@echo "  fmt                - go fmt all packages"
	@echo "  vet                - go vet all packages"
	@echo "  tidy               - go mod tidy"
	@echo "  clean              - Remove bin/ directory"
