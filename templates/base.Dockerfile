FROM {{ default "golang:1.19.1-alpine3.16" .BuilderImage }} AS builder

# Build the binary
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 go build -a -o {{ .AppName }} {{ default "." .MainPackage }}

FROM  {{ default "alpine:3.16.2" .RuntimeImage }} AS runner
WORKDIR /app
COPY --from=builder /workspace .
CMD ["/app/{{ .AppName }}"]