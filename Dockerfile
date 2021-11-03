FROM centos:7
COPY /envelop-rain /root/server
EXPOSE 8080
CMD /root/server