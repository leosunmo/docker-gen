EXTRA_RUN_ARGS?=


build:
	CGO_ENABLED=0 go build  -ldflags "-s -w -X main.BuildVersion=dev -X main.GitCommit=$$(git rev-parse --short HEAD)" -a -o ./bin/docker-gen cmd/gen/gen.go

release: check
	CGO_ENABLED=0 go build  -ldflags "-s -w -X main.BuildVersion=$$(tagver) -X main.GitCommit=$$(git rev-parse --short HEAD)" -a -o ./bin/docker-gen cmd/gen/gen.go

docker: check
	docker build -f gen.Dockerfile --build-arg GIT_COMMIT="$$(git rev-parse --short HEAD)" --build-arg GEN_VERSION="$$(tagver)" -t docker-gen:$$(tagver) .

check:
	@type tagver >/dev/null 2>&1 || (echo "tagver not in PATH, run 'go install github.com/leosunmo/tagver@latest' to install"; exit 1)

.PHONY: build release docker check