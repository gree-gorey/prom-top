FROM golang:1.10.1
RUN go get -d -v github.com/gree-gorey/prom-top/cmd/prom-top
WORKDIR /go/src/github.com/gree-gorey/prom-top/cmd/prom-top
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prom-top .

FROM alpine:3.7
WORKDIR /root/
COPY --from=0 /go/src/github.com/gree-gorey/prom-top/cmd/prom-top/prom-top .
CMD ["./prom-top"]
