# Port Checker with Docker

Scratch based Docker container to test if a port works using a Golang server

[![Build Status](https://travis-ci.org/qdm12/port-checker.svg?branch=master)](https://travis-ci.org/qdm12/port-checker)
[![Docker Build Status](https://img.shields.io/docker/build/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Automated](https://img.shields.io/docker/automated/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)

[![](https://images.microbadger.com/badges/image/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)
[![](https://images.microbadger.com/badges/version/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)

| Download size | Image size | RAM usage | CPU usage |
| --- | --- | --- | --- |
| 2.7MB | 6.97MB | 5.3MB | Very low |

Based on:

- Scratch with the Golang binary
- Ca-Certificates

## Setup

To test port 2345 simply run on the server:

```bash
docker run -it --rm -p 2345:80 qmcgaw/port-checker
```

With a client, access [http://localhost:2345](http://localhost:2345)

## More information

Displays the following information (through HTTP):

- Client IP
- Client Public IP (if client IP is private)
- Client Location
- Client ISP
- Browser and version
- Device type
- OS and version

It uses [https://ipinfo.com](https://ipinfo.com) to obtain extra information about your IP address.

## TO DOs

- [ ] Add CSS to the HTML template
- [ ] Add Google Maps to show the location