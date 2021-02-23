ARG ALPINE_VERSION=3.13
ARG GO_VERSION=1.16
ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add tzdata

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
ENV CGO_ENABLED=0
RUN apk --update add git
WORKDIR /tmp/gobuild
# Copy repository code and install Go dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY index.html index.html
COPY main.go main.go
COPY internal/ ./internal/

FROM --platform=$BUILDPLATFORM base AS test
ENV CGO_ENABLED=1
# g++ is installed for the -race detector in go test
RUN apk --update add g++

FROM --platform=$BUILDPLATFORM base AS lint
ARG GOLANGCI_LINT_VERSION=v1.37.1
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=$BUILDPLATFORM base AS tidy
RUN git init && \
    git config user.email ci@localhost && \
    git config user.name ci && \
    git add -A && git commit -m ci && \
    sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod

FROM --platform=$BUILDPLATFORM base AS build
COPY --from=qmcgaw/xcputranslate:v0.4.0 /xcputranslate /usr/local/bin/xcputranslate
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
RUN GOARCH="$(xcputranslate -targetplatform ${TARGETPLATFORM} -field arch)" \
    GOARM="$(xcputranslate -targetplatform ${TARGETPLATFORM} -field arm)" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o app main.go

FROM scratch
ARG BUILD_DATE
ARG COMMIT
ARG VERSION
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/port-checker/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.title="port-checker" \
    org.opencontainers.image.description="3MB container to check a port works with a Golang server"
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
EXPOSE 8000
ENTRYPOINT ["/port-checker"]
CMD ["-healthserver=true"]
ENV TZ=America/Montreal \
    LISTENING_PORT=8000 \
    ROOT_URL=/
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/port-checker","healthcheck"]
USER 1000
COPY --from=build --chown=1000 /tmp/gobuild/app /port-checker
