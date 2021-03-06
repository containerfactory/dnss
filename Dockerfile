# STEP 1 build executable binary
FROM golang:alpine as builder
MAINTAINER Container Factory <containerfactory@cloudno.de>
# Install git
RUN apk update && apk add git
# Create appuser
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/blitiri.com.ar/go/dnss/
WORKDIR $GOPATH/src/blitiri.com.ar/go/dnss
#get dependancies
RUN go get -d -v
#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '--extldflags "static"' -o /go/bin/dnss
# STEP 2 build a small image
# start from scratch
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/dnss /go/bin/dnss
USER appuser
ENTRYPOINT ["/go/bin/dnss"]
