# syntax=docker/dockerfile:1.4
# Build the manager binary
FROM golang:1.17 as builder
WORKDIR /workspace
COPY . .
RUN --mount=type=cache,id=ztp-dhcp-golang-cache,target=/root/.cache/go-build,sharing=locked,rw CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ztp-dhcp main.go



# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM gcr.io/distroless/static:nonroot
FROM alpine:latest
WORKDIR /
COPY --from=builder /workspace/ztp-dhcp .
USER 65532:65532
ENTRYPOINT ["/ztp-dhcp"]