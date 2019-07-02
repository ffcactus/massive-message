FROM alpine

RUN mkdir -p /opt/message
RUN mkdir -p /opt/message/log
RUN touch /opt/message/log/promise.log
WORKDIR /opt/message
COPY fs ./
RUN chmod a+x ./receiver ./wait-for.sh

EXPOSE 162/udp
