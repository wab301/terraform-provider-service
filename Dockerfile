FROM golang:1.18.3-alpine AS builder

COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

COPY . /app/
RUN go build .

ENV GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct

FROM alpine:3.16

COPY --from=builder /app/terraform-provider-service /usr/local/bin/

ARG ARG_HASH=''

WORKDIR /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/terraform-provider-service" ]