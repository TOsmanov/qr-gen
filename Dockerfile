FROM mirror.gcr.io/golang:1.21 AS builder
COPY ./ /usr/local/go/src/qr_gen/
WORKDIR /usr/local/go/src/qr_gen
RUN go clean --modcache && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
-mod=readonly -o qrgen /usr/local/go/src/qr_gen/server/main.go
FROM scratch
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/go/src/qr_gen/qrgen /app/
CMD ["/app/qrgen"]