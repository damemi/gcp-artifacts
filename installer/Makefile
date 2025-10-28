# Variables
ORG ?= gcr.io/odigos-public
IMAGE_NAME = odigos
TAG ?= latest
PLATFORMS = linux/amd64,linux/arm64

# Full image reference
IMAGE = $(ORG)/$(IMAGE_NAME):$(TAG)

.PHONY: push-image
push-image:
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(IMAGE) \
		--push \
		.

.PHONY: build-image
build-image:
	docker buildx build \
		--platform linux/amd64 \
		--tag $(IMAGE) \
		--load \
		.

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  push-image   - Build multi-arch image and push to registry"
	@echo "  build-image  - Build multi-arch image locally (without push)"
	@echo ""
	@echo "Variables:"
	@echo "  ORG          - Registry organization (default: gcr.io/odigos)"
	@echo "  TAG          - Image tag (default: latest)"
	@echo "  PLATFORMS    - Target platforms (default: linux/amd64,linux/arm64)"
	@echo ""
	@echo "Example usage:"
	@echo "  make push-image"
	@echo "  make push-image ORG=gcr.io/myorg TAG=v1.0.0"

