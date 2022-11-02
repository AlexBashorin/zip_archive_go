FROM alpine:3.16 as root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /go-web/zipit/main

FROM golang:1.17 as builder
WORKDIR /zipit
COPY --from=root-certs /etc/ssl/certs/ca-certs.crt /etc/ssl/certs
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./zipit-stats ./app/./...

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certs.crt /etc/ssl/certs
COPY --chown=1001:1001 --from=builder /zipit/zipit-stats /zipit-stats
USER app
ENTRYPOINT [ "/zipit-stats" ]