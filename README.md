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

To test port 1234, use:

```bash
docker run -d -p 1234:8000/tcp qmcgaw/port-checker
```

or use [docker-compose.yml](https://github.com/qdm12/port-checker/blob/master/docker-compose.yml) with:

```bash
docker-compose up -d
```

With a client, access [http://localhost:1234](http://localhost:1234)

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

## TO DOs

- [ ] Use GeoLite database and Google Maps to show the location
- [ ] Add CSS to the HTML template
- [ ] Precise external mapped port to check it can access itself at start
- [ ] Unit testing
- [ ] Notifications (Pushbullet, email, etc. ?)
- [ ] UDP port check, see [this](https://ops.tips/blog/udp-client-and-server-in-go/)
