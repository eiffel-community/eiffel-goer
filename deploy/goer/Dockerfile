FROM golang:1.24.2-alpine AS build

RUN apk update && apk add --no-cache git make

WORKDIR /tmp/goer
COPY . .
RUN make build

FROM alpine:3.22.2
ENTRYPOINT ["/app/goer"]
COPY --from=build /tmp/goer/bin/goer /app/goer

ARG version
LABEL com.datadoghq.tags.version="${version}" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.version="${version}" \
      org.opencontainers.image.source="https://github.com/eiffel-community/eiffel-goer"
