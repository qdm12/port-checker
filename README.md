# Port Checker with Docker

*3MB container to check a TCP port works with a Golang HTTP server*

<a href="https://github.com/qdm12/port-checker">
  <img src="title.svg" width="300px" height="200px">
</a>

[![Docker Build Status](https://img.shields.io/docker/build/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/port-checker.svg)](https://github.com/qdm12/port-checker/issues)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)
[![Docker Automated](https://img.shields.io/docker/automated/qmcgaw/port-checker.svg)](https://hub.docker.com/r/qmcgaw/port-checker)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)

| Image size | RAM usage | CPU usage |
| --- | --- | --- |
| 2.76MB | 8MB | Very low |

## Setup

1. <details><summary>CLICK IF YOU HAVE AN ARM DEVICE</summary><p>

    - If you have a ARM 32 bit v6 architecture

        ```sh
        docker build -t qmcgaw/port-checker \
        --build-arg BASE_IMAGE_BUILDER=arm32v6/golang \
        --build-arg GOARCH=arm \
        --build-arg GOARM=6 \
        https://github.com/qdm12/port-checker.git
        ```

    - If you have a ARM 32 bit v7 architecture

        ```sh
        docker build -t qmcgaw/port-checker \
        --build-arg BASE_IMAGE_BUILDER=arm32v7/golang \
        --build-arg GOARCH=arm \
        --build-arg GOARM=7 \
        https://github.com/qdm12/port-checker.git
        ```

    - If you have a ARM 64 bit v8 architecture

        ```sh
        docker build -t qmcgaw/port-checker \
        --build-arg BASE_IMAGE_BUILDER=arm64v8/golang \
        --build-arg GOARCH=arm64 \
        https://github.com/qdm12/port-checker.git
        ```

    </p></details>

1. To test port 1234, use:

    ```bash
    docker run -it --rm -p 1234:8000/tcp qmcgaw/port-checker
    ```

    To test port 1234 internally, use:

    ```bash
    docker run -it --rm -e PORT=1234 qmcgaw/port-checker
    ```

1. With a client, access [http://localhost:1234](http://localhost:1234).
You can also port forward with your router to test it is accessible remotely.

## Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `LOGGING` | `json` | `json`, `human` | Logging format |
| `NODEID` | `0` | Any integer | Instance ID for distributed systems |
| `PORT` | `8000` | `1025` to `65535` | TCP port to listen on internally |
| `ROOTURL` | `/` | URL path string | Used if it is running behind a proxy for example |

## More information

Displays the following information (through HTTP):

- Client IP (public or private)
- Browser and version
- Device type
- OS and version

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
