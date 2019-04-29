ARG ALPINE_VERSION=3.9
ARG GO_VERSION=1.12.4

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG BINCOMPRESS
RUN apk --update add git build-base upx
RUN go get -u -v golang.org/x/vgo
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY pkg ./pkg
# RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o app .
RUN [ "${BINCOMPRESS}" == "" ] || (upx -v --best --ultra-brute --overlay=strip app && upx -t app)

FROM scratch
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.schema-version="1.0.0-rc1" \
      maintainer="quentin.mcgaw@gmail.com" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/qdm12/port-checker" \
      org.label-schema.url="https://github.com/qdm12/port-checker" \
      org.label-schema.vcs-description="3MB container to check a port works with a Golang server" \
      org.label-schema.vcs-usage="https://github.com/qdm12/port-checker/blob/master/README.md#setup" \
      org.label-schema.docker.cmd="docker run -d -p 8000:8000/tcp qmcgaw/port-checker" \
      org.label-schema.docker.cmd.devel="docker run -it --rm -p 8000:8000/tcp qmcgaw/port-checker" \
      org.label-schema.docker.params="PORT=1 to 65535 internal listening port" \
      org.label-schema.version="" \
      image-size="2.75MB" \
      ram-usage="8MB" \
      cpu-usage="Very low"
EXPOSE 8000
ENTRYPOINT ["/port-checker"]
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/port-checker","healthcheck"]
USER 1000
COPY index.html /index.html
COPY --from=builder --chown=1000 /tmp/gobuild/app /port-checker