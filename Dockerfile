FROM golang as builder

ENV CGO_ENABLED=0
COPY . /app
WORKDIR /app
RUN go build -x -ldflags="-s -w -v -linkmode=internal -extldflags='-static'"

FROM scratch

COPY --from=builder /app/frosh-api /
EXPOSE 8080

CMD ["/frosh-api"]