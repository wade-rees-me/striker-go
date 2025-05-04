# Directories
SRC_DIR = .
BIN_DIR = bin
TARGET = $(BIN_DIR)/striker-go

# Striker Home
STRIKER = ${HOME}/Striker
LOG_DIR = $(STRIKER)/Simulations/$(shell date +%Y)/$(shell date +%m)/$(shell date +%d)
LOG = $(LOG_DIR)/$(notdir $(TARGET))-$(shell date +%H%M%S).log

# Default goal
.DEFAULT_GOAL := help

.PHONY: all run clean help install lint build run-all run-single-deck run-double-deck run-six-shoe

# Build the Go binary
all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(TARGET) ./$(SRC_DIR)

# Run with default or overridden parameters
run:
	@mkdir -p $(LOG_DIR)
	@echo "Running: --$(STRATEGY) --$(DECKS) -number-of-hands $(HANDS) --number-of-threads $(THREADS)"
	clear
	$(TARGET) --$(STRATEGY) --$(DECKS) --number-of-hands $(HANDS) --number-of-threads $(THREADS) | tee $(LOG)
	sleep 3

# Clean binary
clean:
	rm -rf $(BIN_DIR)

# Lint Go code
lint:
	golangci-lint run ./...

# Install the binary to STRIKER/bin
install:
	@mkdir -p $(STRIKER)/bin
	cp -f $(TARGET) $(STRIKER)/bin

include Makefile.run

