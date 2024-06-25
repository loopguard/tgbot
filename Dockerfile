FROM alpine:latest

WORKDIR ${GOPATH}/src/github.com/loopguard/tgbot/
COPY . ${GOPATH}/src/github.com/loopguard/tgbot/

RUN go build -o ${GOPATH}/bin/service-entrypoint ./cmd