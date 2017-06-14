FROM alpine:3.5

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/bin
COPY ./metrics-collector $GOPATH/bin

ENTRYPOINT ["metrics-collector"]
