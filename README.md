# Port Checker with Docker

*12.8MB container to check a TCP port works with a Golang HTTP server*

<a href="https://github.com/qdm12/port-checker">
  <img src="title.svg" width="300px" height="200px">
</a>

[![Build status](https://github.com/qdm12/port-checker/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/port-checker/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)

## Features

- HTTP lightweight server responding with information on your client:
  - Client IP (public or private)
  - Browser and version
  - Device type
  - OS and version
- Compatible with amd64, 386, armv6, armv7 and arm64 v8 cpu architectures

## Standalone binary

You can download one of the binaries attached with each Github release. For example:

```sh
wget -qO port-checker https://github.com/qdm12/port-checker/releases/download/v0.1.0/port-checker_0.1.0_linux_amd64
chmod +x port-checker
./port-checker --listening-address=":8000" --root-url="/"
```

## Docker image

To test port 1234, use:

```sh
docker run --rm -p 1234:8000/tcp qmcgaw/port-checker
```

To test port 1234 internally, use:

```bash
docker run --rm -e LISTENING_ADDRESS=":1234" qmcgaw/port-checker
```

💁 There are also image tags, such as `:v0.1.0`, corresponding to each tagged release.

## Options

| Environment variable | Flag | Default | Possible values | Description |
| --- | --- | --- | --- | --- |
| `LISTENING_ADDRESS` | `--listening-address` | `:8000` | Valid listening address | TCP address to listen on internally |
| `ROOT_URL` | `--root-url` | `/` | URL path string | Used if it is running behind a proxy for example |

## TO DOs

- [ ] Use GeoLite database and Google Maps to show the location
- [ ] Add CSS to the HTML template
- [ ] Precise external mapped port to check it can access itself at start
- [ ] Unit testing
- [ ] Notifications (Pushbullet, email, etc. ?)
- [ ] UDP port check, see [this](https://ops.tips/blog/udp-client-and-server-in-go/)
