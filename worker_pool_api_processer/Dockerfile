#FROM golang:1.22 AS builder
#WORKDIR /app
#COPY . .
#RUN go build -o app .
#
#FROM alpine:3.19
#COPY --from=builder /app/app .
#RUN addgroup -S appgroup && adduser -S appuser -G appgroup -u 10001
#USER appuser
#CMD ["./app"]


FROM golang:1.25.3 AS builder
RUN useradd -u 10001 nonroot
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY worker_pool_api_processer/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a \
    -ldflags "-extldflags '-static' -w -s" \
    -o app .

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=nonroot:nonroot /app/app /app
USER nonroot
EXPOSE 8080
CMD ["/app"]