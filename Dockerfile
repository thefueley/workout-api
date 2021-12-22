FROM golang:alpine AS builder
LABEL maintainer="Chris Rivera <crivera@eismos.io>"
LABEL description="Comment API Tutorial on tutorialedge.net"

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o app cmd/server/main.go

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]