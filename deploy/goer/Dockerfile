FROM golang:1.17.3-alpine AS build
WORKDIR /tmp/goer
COPY . .
RUN go build -o ./bin/goer ./cmd/goer

FROM alpine:3.13.5
ENTRYPOINT ["/app/goer"]
COPY --from=build /tmp/goer/bin/goer /app/goer
