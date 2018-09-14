FROM alpine:3.8 as alpine
RUN apk --update --no-cache --progress add ca-certificates

FROM golang:alpine AS builder
RUN apk --update add git build-base
WORKDIR /go/src/port-checker
COPY *.go ./
RUN go get -v ./... && \
    go test -v && \
    # build statically
    CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o port-checker .

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/port-checker/port-checker /port-checker
COPY index.html /index.html
EXPOSE 80
ENTRYPOINT ["/port-checker"]
CMD [""]