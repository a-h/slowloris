FROM golang:1.18

ADD . /app

WORKDIR /app

RUN go mod init main.go ; go get

ENTRYPOINT [ "go", "run", "main.go" ]

CMD [ "-delayBeforeFirstByte", "10s", "-delayBeforeLastByte", "10s" ]