FROM golang:alpine as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN apk update && \
    apk upgrade && \
    apk add git
RUN go get -u github.com/go-redis/redis
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./server && \
    mv ./server /usr/bin/server

FROM debian:latest
RUN apt-get update && apt-get install -y texlive-full
# cleanup package manager
RUN apt-get autoclean && apt-get clean
RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /output
COPY --from=builder /usr/bin/server /usr/bin/server
ENTRYPOINT ["server"]
