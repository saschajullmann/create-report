FROM golang:alpine as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN apk update && \
    apk upgrade && \
    apk add git
RUN go get -u github.com/go-redis/redis
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./client && \
    mv ./client /usr/bin/client

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/bin/client /usr/bin/client
ENTRYPOINT ["client"]
