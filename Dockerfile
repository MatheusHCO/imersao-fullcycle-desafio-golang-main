FROM golang:1.19

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

COPY . /go/src/

RUN apt-get update && \
    apt-get install sqlite3 libsqlite3-dev build-essential -y

CMD ["/go/src/./main"]