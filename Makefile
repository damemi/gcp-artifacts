ORG ?= gcr.io/odigos-public
IMAGE_NAME = odigos/deployer
TAG ?= 1.11.2
PLATFORMS = linux/amd64

# Full image reference
IMAGE = $(ORG)/$(IMAGE_NAME):$(TAG)

.PHONY: push-image
push-image:
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(IMAGE) \
		--push \
		--annotation "index:com.googleapis.cloudmarketplace.product.service.name=services/odigos.endpoints.odigos-public.cloud.goog" \
		.
