FROM golang:alpine3.16 as builder
COPY ./ /metric-exporter-demo
WORKDIR /metric-exporter-demo
RUN go build -o ./bin/exporter main.go



FROM alpine
COPY --from=builder /metric-exporter-demo/bin/exporter /metric-exporter
ENTRYPOINT bash -c "/metric-exporter"
