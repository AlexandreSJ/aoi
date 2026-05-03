.PHONY: build run clean

BINARY   := aoi
BUILD_DIR := ./build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/aoi/

run: build
	$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)
