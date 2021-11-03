# FROM golang:alpine AS builder
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn,direct
# ENV CGO_ENABLED 0
# COPY . /root
# # RUN  apk --update add git tzdata
# WORKDIR /root
# RUN go build -o /root/app

FROM cr-cn-beijing.volces.com/group6/centos:7
# COPY --from=builder /root/app /root/server
# COPY --from=builder /root/configs/ /root/configs/
COPY envelop-rain /root/server
COPY ./configs/ /root/configs/
EXPOSE 8080
CMD /root/server