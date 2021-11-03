FROM golang:alpine AS builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
# ENV CGO_ENABLED 0
# COPY . /root
# RUN  apk --update add git tzdata
# WORKDIR /root
RUN go build

FROM ubuntu
COPY --from=builder envelop-rain /root/server
COPY --from=builder /configs/ /root/configs/
WORKDIR /root
# COPY envelop-rain /root/server
# COPY configs/ /root/configs/
EXPOSE 8080
CMD ./server