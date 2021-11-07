# FROM golang:alpine AS builder
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn,direct
# ENV CGO_ENABLED 0
# COPY . /root
# # RUN  apk --update add git tzdata
# WORKDIR /root
# RUN go build -o /root/app

# webhook test
FROM cr-cn-beijing.volces.com/group7/centos:7
# COPY --from=builder /root/app /root/server
# COPY --from=builder /root/configs/ /root/configs/
WORKDIR /root
COPY envelop-rain ./server
COPY ./configs/ ./configs/
EXPOSE 8080
CMD /root/server