FROM alpine as ssl

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM golang as builder

ENV CGO_ENABLED=0
COPY . /app
WORKDIR /app
RUN go build -x -ldflags="-s -w -v -linkmode=internal -extldflags='-static'"

FROM scratch

COPY --from=builder /app/frosh-api /
COPY --from=ssl /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

CMD ["/frosh-api"]