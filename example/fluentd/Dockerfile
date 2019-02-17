FROM golang:1.11.5 AS builder

ENV GOPATH /go

RUN go get -d github.com/abicky/ecswrap \
  && cd $GOPATH/src/github.com/abicky/ecswrap \
  && CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s"

FROM fluent/fluentd:v1.3.3-1.0

COPY --from=builder /go/src/github.com/abicky/ecswrap/ecswrap /usr/local/bin/
ENTRYPOINT ["tini", "--", "ecswrap", "-v", "--", "/bin/entrypoint.sh"]
CMD ["fluentd"]
