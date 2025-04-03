FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24 AS builder

ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ENV VERSION=$VERSION

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-X github.com/clambin/homewizard-exporter/internal/collector.version=$VERSION" \
    -o homewizard-exporter \
    .

FROM alpine

RUN apk add --no-cache tzdata

WORKDIR /app
COPY --from=builder /app/homewizard-exporter /app/homewizard-exporter

RUN /usr/sbin/addgroup app
RUN /usr/sbin/adduser app -G app -D
USER app

ENTRYPOINT ["/app/homewizard-exporter"]
CMD []
