FROM busybox
# RUN go get github.com/segmentio/throttle
# RUN go get github.com/segmentio/json_to_nsq
COPY simple-echo-microservice.sh /entrypoint.sh
COPY simple-echo-microservice.sh /service-a.sh
# Dont run commands - use exec instead
# CMD ["/bin/sh", "/entrypoint.sh"]