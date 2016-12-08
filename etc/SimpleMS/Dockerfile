FROM busybox
# RUN go get github.com/segmentio/json_to_nsq
COPY simple-echo-microservice.sh /entrypoint.sh
COPY simple-echo-microservice.sh /service-a.sh
CMD ["/bin/sh", "/entrypoint.sh"]
# CMD ["/entrypoint.sh"]