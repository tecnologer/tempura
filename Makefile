.SILENT:
.PHONY: *

OUTPUT_DIR=./bin
RELEASE_BRANCH=main
CURRENT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Get the latest tag matching the pattern vX.Y.Z
LATEST_VERSION=$(shell git tag --list 'v*' --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -n 1)

# Check if RELEASE_VERSION is provided, else increment the patch version of the latest tag
#
# Example: If the RELEASE_VERSION is not provided and the latest tag is v1.2.3, the next version will be v1.2.4
NEXT_VERSION=$(if $(RELEASE_VERSION),$(RELEASE_VERSION),$(shell echo $(LATEST_VERSION) | sed 's/^v//' | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}'))

migrator:
	go build -ldflags "-X 'main.version=$(NEXT_VERSION)'" -o $(OUTPUT_DIR)/migrator ./cmd/migrator/main.go
