ARG ALPINE_VERSION=3.11
ARG GO_VERSION=1.14

FROM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add ca-certificates tzdata

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV CGO_ENABLED=0
RUN apk --update add git
ARG GOLANGCI_LINT_VERSION=v1.27.0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}
WORKDIR /tmp/gobuild
COPY .golangci.yml .
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY main.go .
COPY pkg ./pkg
# RUN go test -v -race ./...
RUN go build -a -ldflags="-s -w" -o app main.go

FROM scratch
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/port-checker/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.title="port-checker" \
    org.opencontainers.image.description="3MB container to check a port works with a Golang server"
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
EXPOSE 8000
ENTRYPOINT ["/port-checker"]
ENV TZ=America/Montreal \
    LISTENING_PORT=8000 \
    ROOT_URL=/
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/port-checker","healthcheck"]
USER 1000
COPY --chown=1000 index.html /index.html
COPY --from=builder --chown=1000 /tmp/gobuild/app /port-checker