FROM cr-cn-beijing.volces.com/group7/centos:7
WORKDIR /root
COPY envelop-rain ./server
COPY ./configs/ ./configs/
EXPOSE 8080
CMD /root/server