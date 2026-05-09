.PHONY: build run clean version patch minor major release

BINARY    := aoi
BUILD_DIR := ./build
LDFLAGS   := -ldflags "-s -w"

build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/$(BINARY)

run: build
	$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

version:
	@bash -c 'TAG=$$(git tag --sort=-version:refname | head -1); echo "$${TAG:-no tags yet}"'

patch:
	@bash -c 'TAG=$$(git tag --sort=-version:refname | head -1); CURRENT=$${TAG:-v0.0.0}; V=$${CURRENT#v}; IFS=. read MAJOR MINOR PATCH <<< "$$V"; NEW="v$$MAJOR.$$MINOR.$$((PATCH+1))"; git tag $$NEW && echo "Created tag: $$NEW"'

minor:
	@bash -c 'TAG=$$(git tag --sort=-version:refname | head -1); CURRENT=$${TAG:-v0.0.0}; V=$${CURRENT#v}; IFS=. read MAJOR MINOR PATCH <<< "$$V"; NEW="v$$MAJOR.$$((MINOR+1)).0"; git tag $$NEW && echo "Created tag: $$NEW"'

major:
	@bash -c 'TAG=$$(git tag --sort=-version:refname | head -1); CURRENT=$${TAG:-v0.0.0}; V=$${CURRENT#v}; IFS=. read MAJOR MINOR PATCH <<< "$$V"; NEW="v$$((MAJOR+1)).0.0"; git tag $$NEW && echo "Created tag: $$NEW"'

release:
	git push origin main --tags
