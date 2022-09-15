FROM golang:1.19.1-alpine3.16 AS builder

ARG GEN_VERSION="dev"
ARG GIT_COMMIT="unknown"

# Build the binary
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags="-X main.GitCommit=${GIT_COMMIT} -X main.Version=${GEN_VERSION}" -o docker-gen cmd/gen/gen.go

FROM  alpine:3.16.2 AS runner
WORKDIR /app
COPY --from=builder /workspace .
ENTRYPOINT ["/app/docker-gen"]
