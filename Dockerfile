FROM golang:alpine AS builder
WORKDIR /go/src/port-checker
COPY *.go ./
RUN apk --update add git build-base
RUN go get -v ./... && \
    go test -v && \
    # build statically
    CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o port-checker .

FROM alpine:3.8 as alpine
RUN apk --update --no-cache --progress add ca-certificates

FROM scratch
COPY --from=builder /go/src/port-checker/port-checker /port-checker
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 80
ENTRYPOINT ["/port-checker"]
CMD [""]
