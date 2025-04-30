FROM golang:1.23.6-alpine AS builder

WORKDIR /app/
ENV PACKAGES="curl build-base git bash file linux-headers eudev-dev"
RUN apk add --no-cache $PACKAGES

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make build

FROM alpine:latest

COPY --from=builder /app/build/sunrised /usr/bin/sunrised
EXPOSE 26656 26657 1317 9090

CMD ["sleep", "infinity"]
