# builder image
FROM golang:1.21-alpine3.18 as builder
ARG BUILD_VERSION
LABEL org.opencontainers.image.source="https://github.com/nilic/hntop-cli"
RUN mkdir /build
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$BUILD_VERSION" -a -o hntop ./cmd/hntop

# generate clean, final image for end users
FROM alpine:3.18
COPY --from=builder /build/hntop .

# executable
ENTRYPOINT [ "./hntop" ]