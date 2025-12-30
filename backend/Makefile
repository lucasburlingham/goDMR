# Variables
OUTPUT_DIR=./bin
BINARY_x86=dmr-api
BINARY_arm=dmr-api-arm
SRC=main.go
INSTALL_DIR=/opt/dmr-backend
CONFIG=config.ini
USER=dmr
GOARM=6

# Build for local or Pi Zero
all: build

# Initialize module (if go.mod does not exist)
init:
	@if [ ! -f go.mod ]; then \
		echo "Initializing Go module..."; \
		go mod init dmr-backend; \
	fi
	go get gopkg.in/ini.v1@v1.66.4

# Build binary for local machine
build: init
	go build -o $(OUTPUT_DIR)/$(BINARY_x86) $(SRC)

# Cross-compile for Pi Zero
build-pi: init
	GOOS=linux GOARCH=arm GOARM=$(GOARM) go build -o $(OUTPUT_DIR)/$(BINARY_arm) $(SRC)
# Install binary and config
install: build
	sudo mkdir -p $(INSTALL_DIR)
# 	Find correct binary based on architecture
	if [ "$(shell uname -m)" = "armv6l" ]; then \
		sudo cp $(OUTPUT_DIR)/$(BINARY_arm) $(INSTALL_DIR)/dmr-api; \
	else \
		sudo cp $(OUTPUT_DIR)/$(BINARY_x86) $(INSTALL_DIR)/dmr-api; \
	fi
	sudo cp $(CONFIG) $(INSTALL_DIR)/
	sudo chown -R $(USER):$(USER) $(INSTALL_DIR)
# Reload systemd to pick up service file
systemd-reload:
	sudo systemctl daemon-reload
	sudo systemctl enable dmr-backend.service
	sudo systemctl restart dmr-backend.service

# Clean build artifacts
clean:
	rm -f $(OUTPUT_DIR)/$(BINARY)
	rm -f go.sum

.PHONY: all init build build-pi install systemd-reload clean
