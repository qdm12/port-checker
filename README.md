# Port Checker with Docker

1.89MB container to check a port works with a Golang server

[![Build Status](https://travis-ci.org/qdm12/port-checker.svg?branch=master)](https://travis-ci.org/qdm12/port-checker)
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
| 3.23MB | 8MB | Very low |

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

## More information

Displays the following information (through HTTP):

- Client IP (public or private)
- Browser and version
- Device type
- OS and version

## TO DOs

- [ ] Emojis
- [ ] Precise port to check it can access itself at start
- [ ] Add CSS to the HTML template
- [ ] Use GeoLite database
- [ ] Add Google Maps to show the location
- [ ] Unit testing and code refactoring
- [ ] Notifications (Pushbullet, email, etc. ?)
