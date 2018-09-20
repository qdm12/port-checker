FROM alpine:3.8 as alpine
RUN apk --update --no-cache --progress add ca-certificates

FROM golang:alpine AS builder
RUN apk --update add git build-base upx
WORKDIR /go/src/port-checker
COPY *.go ./
RUN go get -v ./... && \
    go test -v && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o updater . && \
    upx -v --best --ultra-brute --overlay=strip updater && \
    upx -t updater

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/port-checker/port-checker /port-checker
EXPOSE 80
ENTRYPOINT ["/port-checker"]
COPY index.html /index.html