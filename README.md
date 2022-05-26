Run a web server which waits a configurable amount of time before responding. See also https://adrianhesketh.com/2016/12/03/testing-slow-http-responses/

# Build
```
docker build -t slowloris .
```

# Run
```
docker run -it -p8080:8080 --rm slowloris
```

# Use
```
curl -I localhost:8080/foo
```