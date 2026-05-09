.PHONY: build run clean patch minor major

BINARY   := aoi
BUILD_DIR := ./build
STYLES_GO  := internal/ui/styles.go

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/$(BINARY)

run: build
	$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

# Version management
CURRENT_VERSION := $(shell grep 'const footerVersion' $(STYLES_GO) | sed 's/.*"\(.*\)"/\1/')

version:
	@echo "Current version: $(CURRENT_VERSION)"

patch: version
	@bash -c 'VERSION_NUMBERS="$$(echo "$(CURRENT_VERSION)" | sed "s/^v//")"; MAJOR=$$(echo $$VERSION_NUMBERS | cut -d. -f1); MINOR=$$(echo $$VERSION_NUMBERS | cut -d. -f2); PATCH=$$(echo $$VERSION_NUMBERS | cut -d. -f3); NEW_PATCH=$$((PATCH + 1)); NEW_VERSION="v$$MAJOR.$$MINOR.$$NEW_PATCH"; sed -i "s/const footerVersion = \".*\"/const footerVersion = \"$$NEW_VERSION\"/" $(STYLES_GO); echo "Updated to: $$NEW_VERSION"'

minor: version
	@bash -c 'VERSION_NUMBERS="$$(echo "$(CURRENT_VERSION)" | sed "s/^v//")"; MAJOR=$$(echo $$VERSION_NUMBERS | cut -d. -f1); MINOR=$$(echo $$VERSION_NUMBERS | cut -d. -f2); NEW_MINOR=$$((MINOR + 1)); NEW_VERSION="v$$MAJOR.$$NEW_MINOR.0"; sed -i "s/const footerVersion = \".*\"/const footerVersion = \"$$NEW_VERSION\"/" $(STYLES_GO); echo "Updated to: $$NEW_VERSION"'

major: version
	@bash -c 'VERSION_NUMBERS="$$(echo "$(CURRENT_VERSION)" | sed "s/^v//")"; MAJOR=$$(echo $$VERSION_NUMBERS | cut -d. -f1); NEW_MAJOR=$$((MAJOR + 1)); NEW_VERSION="v$$NEW_MAJOR.0.0"; sed -i "s/const footerVersion = \".*\"/const footerVersion = \"$$NEW_VERSION\"/" $(STYLES_GO); echo "Updated to: $$NEW_VERSION"'

release:
	@bash git tag $(CURRENT_VERSION)
	@echo New tag created. Push using --tags