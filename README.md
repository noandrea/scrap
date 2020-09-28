# Scrap

A simple tool to scrap movie data from content providers

## Requirements 

Scrap requires a headless chrome to correctly render and parse html pages.
Use this command from the terminal to run chrome in headless mode:

```sh
chrome --headless --disable-gpu --remote-debugging-port=9222
```

## Usage

Scarp can be used as a REST API service, command line client or library. 

#### Providers
Currently scrap supports only Amazon Prime (DE)

#### Installation 

run `go get github.com/noandrea/scrap` or download a linux binary from the [release page](https://github.com/noandrea/scrap/releases) 

### Command line client 

The command line client provides the commands to:
- `inspect`: inspect an id from a provider
- `serve`: run the REST API server

#### Inspect

The `inspect` command is intended mostly for debug purposes, here an example how to use it:

[![asciicast](https://asciinema.org/a/361847.svg)](https://asciinema.org/a/361847)

> 💡 use the `-q` option to use it with `jq`, ex: `scrap inspect B00JKEJ4TA -q | jq ".title"` 

#### Serve

The `serve` command runs a web server that exposes the REST API endpoints. 

[![asciicast](https://asciinema.org/a/361848.svg)](https://asciinema.org/a/361848)

The web server exposes the following rest API endpoints:
- `/movie/amazon/{amazon_id}`: scrape the data from amazon
- `/status`: a monitoring endpoint

> TODO specify replies and status codes

**Configuration** 

An external configuration file can be provided to tweak the behavior of the serve:

```yaml
listen_address: :8080 # server listen address
chrome_address: 127.0.0.1:9222 # address of chrome headless 
scrap_region: de # the region to use (for amazon is the domain)
cache_enabled: false # enable cache for movies (NOT YET IMPLEMENTED)
cache_max_size: 5000 # maximum number of cached movies (NOT YET IMPLEMENTED)
cache_lifetime: 86400 # lifetime of a cached movie (in seconds) (NOT YET IMPLEMENTED)
```

### Programmatic usage

> TODO add programmatic usage

## Docker

The app can be built using docker with the command:

```sh
make docker-build
```

The docker build is optimized for space and is built on top of `scratch`

After building, you can run the docker image with 

```
docker run -p 8080:8080 scrap:latest
```

:warning: chrome headless is not bundled with the image!!

## Known Issues

- The FQL query for amazon prime is not robust enough and it fails for some  IDs
- The docker image and the examples lack the configuration for chrome headless


## Examples

### Systemd

This assumes that scrap is installed in `/usr/bin/scrap`

```systemd
[Unit]
Description=scrap

[Service]
User=www-data
Group=www-data
Type=simple
ExecStart=/usr/bin/scrap start --debug
Restart=always

[Install]
WantedBy=multi-user.target
```                           

### Docker compose

`docker-compose.yaml` example

```yaml
version: '3'
services:
  scrap:
    container_name: scrap
    image: scrap:latest
    ports:
    - 8080:8080
    environment: 
    - SCRAP_CHROME_ADDRESS=http://chrome:9222
  chrome:
    container_name: scrap_chrome
    image: zenika/alpine-chrome
    expose:
    - 9222
    command: 
         [chromium-browser, "--headless", "--disable-gpu", "--no-sandbox", "--remote-debugging-address=0.0.0.0", "--remote-debugging-port=9222"]
```


### K8s

Kubernetes configuration example:

```yaml
---
# Deployment
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: scrap
  name: scrap
spec:
  replicas: 1
  revisionHistoryLimit: 3
  selector:
    matchLabels:
      app: scrap
  template:
    metadata:
      labels:
        app: scrap
    spec:
      containers:
      - env:
        image: scrap:latest
        imagePullPolicy: Always
        name: scrap
        ports:
        - name: http
          containerPort: 8080
        livenessProbe:
          httpGet:
            path: /status
            port: 8080
---
# Service
# the service for the above deployment
apiVersion: v1
kind: Service
metadata:
  name: scrap-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  selector:
    app: scrap

```
