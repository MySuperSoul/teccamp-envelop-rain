FROM golang:alpine AS builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
COPY . /envelop-rain
RUN  apk --update add git tzdata
WORKDIR /envelop-rain
RUN go build -o /envelop-rain/app

FROM centos:7
COPY --from=builder /envelop-rain/app /root/server
COPY --from=builder /envelop-rain/configs/ /root/configs/
WORKDIR /root
EXPOSE 8080
CMD /root/server