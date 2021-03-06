# Port Checker with Docker

*12.8MB container to check a TCP port works with a Golang HTTP server*

<a href="https://github.com/qdm12/port-checker">
  <img src="title.svg" width="300px" height="200px">
</a>

[![Build status](https://github.com/qdm12/port-checker/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/port-checker/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
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

## Setup

1. To test port 1234, use:

    ```bash
    docker run -it --rm -p 1234:8000/tcp qmcgaw/port-checker
    ```

    To test port 1234 internally, use:

    ```bash
    docker run -it --rm -e LISTENING_PORT=1234 qmcgaw/port-checker
    ```

1. With a client, access [http://localhost:1234](http://localhost:1234).
You can also port forward with your router to test it is accessible remotely.

### Binary

You can also download one of the binaries on the Github releases. For example:

```sh
wget -qO port-checker https://github.com/qdm12/port-checker/releases/download/v0.1.0/port-checker_0.1.0_linux_amd64
chmod +x port-checker
./port-checker
# Usage with
./port-checker -help
```

## Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `LISTENING_PORT` | `8000` | `1025` to `65535` | TCP port to listen on internally |
| `ROOT_URL` | `/` | URL path string | Used if it is running behind a proxy for example |

## Development

### Using VSCode and Docker

1. Install [Docker](https://docs.docker.com/install/)
    - On Windows, share a drive with Docker Desktop and have the project on that partition
1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
1. Your dev environment is ready to go!... and it's running in a container :+1:

## TO DOs

- [ ] Use GeoLite database and Google Maps to show the location
- [ ] Add CSS to the HTML template
- [ ] Precise external mapped port to check it can access itself at start
- [ ] Unit testing
- [ ] Notifications (Pushbullet, email, etc. ?)
- [ ] UDP port check, see [this](https://ops.tips/blog/udp-client-and-server-in-go/)
