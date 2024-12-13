# Build stage
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache gcc musl-dev upx

WORKDIR /build
COPY go.* .
COPY *.go .
COPY templates/ templates/

# Build with size optimizations
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main .
# Compress the binary with UPX
RUN upx --best --lzma main

# Final stage - using scratch (smallest possible base)
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /build/main .
COPY --from=builder /build/templates/ templates/

EXPOSE 8080
CMD ["/app/main"]