#!/usr/bin/env bash
# Builds patched versions of the official Odigos component images that the Marketplace vulnerability
# scan flags, and publishes them like hack/publish-component-images.sh does (annotation + tags).
# Run publish-component-images.sh first; this overwrites the tags of the patched components only.
#
# Usage: hack/patch-component-images.sh 1.39.0 <path to dtrace-injector.jar>
#   The jar replaces /instrumentations/java-ebpf/dtrace-injector.jar in the enterprise images; it must be
#   built from the same java-ebpf version the release ships.
set -euo pipefail

VERSION="${1:?usage: $0 <version> <dtrace-injector.jar>}"
JAR="${2:?usage: $0 <version> <dtrace-injector.jar>}"
VERSION="${VERSION#v}"
TRACK="${VERSION%.*}"
SRC="${SRC:-registry.odigos.io}"
DEST="${DEST:-gcr.io/odigos-public/odigos}"
STAGING="${STAGING:-gcr.io/odigos-public/odigos-staging}"
SERVICE_NAME="services/odigos.endpoints.odigos-public.cloud.goog"
HERE="$(cd "$(dirname "$0")" && pwd)"

ctx="$(mktemp -d)"
trap 'rm -rf "$ctx"' EXIT
cp "$JAR" "$ctx/dtrace-injector.jar"

# component  patch-grpc  patch-java
PATCHES=(
  "odiglet            yes no"
  "agents             no  no"
  "enterprise-odiglet no  yes"
  "enterprise-agents  no  yes"
)

for p in "${PATCHES[@]}"; do
  read -r c grpc java <<<"$p"
  src="${SRC}/odigos-${c}:v${VERSION}"
  base="${src}@$(crane digest "$src")"
  user="$(crane config --platform linux/amd64 "$src" | python3 -c 'import json,sys;print(json.load(sys.stdin)["config"].get("User") or "0")')"
  stage="${STAGING}/odigos-${c}:v${VERSION}"
  echo "==> ${c}: ${base} (grpc=${grpc} java=${java} user=${user})"
  docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --provenance=false --sbom=false \
    --build-arg BASE="$base" \
    --build-arg PATCH_GRPC="$grpc" \
    --build-arg PATCH_JAVA="$java" \
    --build-arg RUNTIME_USER="$user" \
    -f "$HERE/patch-component-images.Dockerfile" \
    -t "${stage}-overlay" --push \
    "$ctx"
  crane flatten "${stage}-overlay" -t "${stage}-flat"
  docker buildx imagetools create \
    --annotation "index:com.googleapis.cloudmarketplace.product.service.name=${SERVICE_NAME}" \
    --tag "${DEST}/odigos-${c}:v${VERSION}" \
    --tag "${DEST}/odigos-${c}:${VERSION}" \
    --tag "${DEST}/odigos-${c}:${TRACK}" \
    "${stage}-flat"
done
