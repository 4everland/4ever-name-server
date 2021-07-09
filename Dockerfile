FROM golang:alpine AS builder
LABEL stage=dnsbuilder

RUN apk update --no-cache
RUN apk add --no-cache git
RUN apk add --no-cache make

WORKDIR /build

ADD . .

RUN chmod 755 /build/build.sh && /build/build.sh

FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /build/coredns/coredns /app/coredns

CMD ["./coredns"]