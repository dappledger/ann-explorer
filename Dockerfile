FROM golang:1.12.1-alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git
ADD . /ann
ENV GOPATH=/gopkg
RUN cd /ann && \
    ./get_pkgs.sh

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /ann/conf /block-browser/conf
COPY --from=builder /ann/static /block-browser/static
COPY --from=builder /ann/views /block-browser/views
COPY --from=builder /ann/block-browser /usr/local/bin
WORKDIR /block-browser

ENTRYPOINT ["block-browser"]
