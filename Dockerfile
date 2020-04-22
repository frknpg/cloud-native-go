FROM golang:1.14.2-alpine
MAINTAINER furkan

ENV SOURCES /go/src/cloud-native-go

COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go install

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT cloud-native-go
