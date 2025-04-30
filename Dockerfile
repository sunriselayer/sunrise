FROM golang:1.23.6-alpine AS builder

WORKDIR /app/
ENV PACKAGES="curl build-base git bash file linux-headers eudev-dev"
RUN apk add --no-cache $PACKAGES

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN addgroup -g 1025 nonroot
RUN adduser -D nonroot -u 1025 -G nonroot

COPY --from=builder /app/build/sunrised /usr/bin/sunrised
EXPOSE 26656 26657 1317 9090
USER nonroot

ENTRYPOINT ["sunrised"]