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

[![](https://images.microbadger.com/badges/image/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)
[![](https://images.microbadger.com/badges/version/qmcgaw/port-checker.svg)](https://microbadger.com/images/qmcgaw/port-checker)

| Download size | Image size | RAM usage | CPU usage |
| --- | --- | --- | --- |
| ?MB | 1.89MB | 7.7MB | Very low |

## Setup

### Check the port

To test port 2345 simply run on the server:

```bash
docker run -it --rm -p 2345:80 qmcgaw/port-checker
```

With a client, access [http://localhost:2345](http://localhost:2345)

### Check the port and the IP address

To test port 2345 simply run on the server:

```bash
docker run -it --rm --net host -e PORT=2345 qmcgaw/port-checker
```

With a client, access [http://localhost:2345](http://localhost:2345) and your IP address will be shown and not the Docker gateway one

You can also port forward with your router to test it is accessible remotely.

There is a *docker-compose.yml* file if you are interested.

## More information

Displays the following information (through HTTP):

- Client IP (public or private)
- Browser and version
- Device type
- OS and version

## TO DOs

- [ ] Emojis
- [ ] Healthcheck
- [ ] Precise port to check it can access itself at start
- [ ] Add CSS to the HTML template
- [ ] Use GeoLite database
- [ ] Add Google Maps to show the location
- [ ] Unit testing and code refactoring
- [ ] Notifications (Pushbullet, email, etc. ?)
