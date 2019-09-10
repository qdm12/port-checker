ARG BASE_IMAGE_BUILDER=golang
ARG ALPINE_VERSION=3.10
ARG GO_VERSION=1.13

FROM ${BASE_IMAGE_BUILDER}:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG GOARCH=amd64
ARG GOARM
ARG BINCOMPRESS
RUN apk --update add git build-base upx
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY main.go .
COPY pkg ./pkg
# RUN go test -v -race ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GOARM=${GOARM} go build -a -ldflags="-s -w" -o app .
RUN [ "${BINCOMPRESS}" == "" ] || (upx -v --lzma --overlay=strip app && upx -t app)

FROM scratch
ARG BUILD_DATE
ARG VCS_REF
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version="" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/port-checker/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.title="port-checker" \
    org.opencontainers.image.description="3MB container to check a port works with a Golang server" \
    image-size="2.76MB" \
    ram-usage="8MB" \
    cpu-usage="Very low"
EXPOSE 8000
ENTRYPOINT ["/port-checker"]
ENV PORT=8000 \
    LOGGING=json \
    NODEID=0 \
    ROOTURL=/
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/port-checker","healthcheck"]
USER 1000
COPY --chown=1000 index.html /index.html
COPY --from=builder --chown=1000 /tmp/gobuild/app /port-checker