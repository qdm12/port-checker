ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23
ARG BUILDPLATFORM=linux/amd64
ARG XCPUTRANSLATE_VERSION=v0.6.0
ARG GOLANGCI_LINT_VERSION=v1.61.0

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate
FROM --platform=${BUILDPLATFORM} qmcgaw/binpot:golangci-lint-${GOLANGCI_LINT_VERSION} AS golangci-lint

FROM --platform=${BUILDPLATFORM} alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add tzdata

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
ENV CGO_ENABLED=0
RUN apk --update add git g++
WORKDIR /tmp/gobuild
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
COPY --from=golangci-lint /bin /go/bin/golangci-lint
# Copy repository code and install Go dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY index.html index.html
COPY main.go main.go
COPY internal/ ./internal/

FROM --platform=${BUILDPLATFORM} base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we installed g++ to support the race detector
ENV CGO_ENABLED=1
ENTRYPOINT go test -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...

FROM --platform=$BUILDPLATFORM base AS lint
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=${BUILDPLATFORM} base AS build
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG CREATED="an unknown date"
ARG COMMIT=unknown
RUN echo ${TARGETPLATFORM} && sleep 10
RUN GOARCH="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -field arch)" \
    GOARM="$(xcputranslate translate -targetplatform ${TARGETPLATFORM} -field arm)" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.created=$CREATED' \
    -X 'main.commit=$COMMIT' \
    " -o app main.go

FROM scratch
ARG CREATED
ARG COMMIT
ARG VERSION
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$CREATED \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/port-checker/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/port-checker" \
    org.opencontainers.image.title="port-checker" \
    org.opencontainers.image.description="3MB container to check a port works with a Golang server"
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
EXPOSE 8000/tcp
ENTRYPOINT ["/port-checker"]
ENV TZ=America/Montreal \
    LISTENING_ADDRESS=:8000 \
    ROOT_URL=/
ARG UID=1000
ARG GID=1000
USER ${UID}:${GID}
COPY --from=build --chown=${UID}:${GID} /tmp/gobuild/app /port-checker
