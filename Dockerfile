FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ADD . /app
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-s -w" -trimpath -o /usr/local/bin/frosh-api

FROM gcr.io/distroless/base

COPY --from=builder /usr/local/bin/frosh-api /usr/local/bin/frosh-api

ENTRYPOINT ["/usr/local/bin/frosh-api"]
