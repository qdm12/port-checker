# Port Checker with Docker

Simple Docker container to test if a port works using a Golang server

Displays the following information (through HTTP):

- Client IP
- Client Public IP (if client IP is private)
- Client Location
- Client ISP
- Browser and version
- Device type
- OS and version

It uses [https://ipinfo.com](https://ipinfo.com) to obtain extra information about your IP address.

## Setup

To test port 2345 simply run:

```bash
docker run -it --rm -p 2345:80 qmcgaw/port-checker
```

## TO DOs

- [ ] Add CSS to the HTML template
- [ ] Add Google Maps to show the location