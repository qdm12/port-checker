FROM golang:alpine AS builder
RUN echo https://dl-3.alpinelinux.org/alpine/v3.8/main > /etc/apk/repositories && \
    apk --update add git build-base && \
    echo https://dl-3.alpinelinux.org/alpine/v3.8/community >> /etc/apk/repositories && \
    apk --update add upx
RUN go get -u -v golang.org/x/vgo
WORKDIR /tmp/gobuild
    
FROM scratch AS final
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="1.89MB container to check a port works with a Golang server" \
      download="1.89MB" \
      size="MB" \
      ram="7.7MB" \
      cpu_usage="Very low" \
      github="https://github.com/qdm12/port-checker"
EXPOSE 80
ENTRYPOINT ["/port-checker"]

FROM builder AS builder2
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o app .
RUN upx -v --best --ultra-brute --overlay=strip app && upx -t app

FROM final
COPY --from=builder2 /tmp/gobuild/app /port-checker
COPY index.html /index.html
