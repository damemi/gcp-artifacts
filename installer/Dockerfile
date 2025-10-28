FROM --platform=$BUILDPLATFORM golang:1.25.3 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -tags=embed_manifests -o installer .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/installer /installer

ENTRYPOINT ["/installer"]


