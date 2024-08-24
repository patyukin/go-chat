FROM golang:1.22.3-alpine3.19 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/chat cmd/chat/main.go

FROM alpine:3.19

RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=builder /app/bin/chat .
ENV YAML_CONFIG_FILE_PATH=config.yaml
COPY migrations migrations
COPY internal/templates internal/templates

CMD ["./chat"]
