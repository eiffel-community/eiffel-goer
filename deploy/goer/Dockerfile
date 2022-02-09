FROM golang:1.17.3-alpine AS build

RUN apk update && apk add --no-cache git make

WORKDIR /tmp/goer
COPY . .
RUN make build

FROM alpine:3.13.5
ENTRYPOINT ["/app/goer"]
COPY --from=build /tmp/goer/bin/goer /app/goer
