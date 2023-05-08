FROM golang:1.19-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GO111MODULE=on go build -o main ./cmd/server/main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
RUN chown -R appuser /app & mkdir -p /app/config & mkdir -p /app/tmp & mkdir -p /app/templates
USER appuser
COPY --from=builder /build/main /app/
COPY config /app/config
COPY templates /app/templates
WORKDIR /app
EXPOSE 8088
CMD ["./main"]
