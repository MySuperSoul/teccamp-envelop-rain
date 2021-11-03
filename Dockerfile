FROM golang:alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go build
EXPOSE 8080
CMD /envelop-rain