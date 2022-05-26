FROM golang:1.18

ADD . /app

WORKDIR /app

RUN go build

ENTRYPOINT [ "/app/slowloris" ]

CMD [ "-delayBeforeFirstByte", "10s", "-delayBeforeLastByte", "10s" ]
