# slorloris

Run a web server which waits a configurable amount of time before responding. See also https://adrianhesketh.com/2016/12/03/testing-slow-http-responses/

## Install

### With Go

* Run `go install github.com/a-h/slowloris@latest`

## Build

### With Go

* Install Go.
* Run `go build`.
* `./slowloris`

### With Docker

```
docker build -t slowloris .
```

```
docker run -it -p8080:8080 --rm slowloris
```

## Usage

```
curl -I localhost:8080/foo
```