FROM centos:7
COPY /envelop-rain /root/server
COPY /configs /root/server
EXPOSE 8080
CMD /root/server