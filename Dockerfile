FROM golang:1.14.2-alpine3.11 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR ./src/github.com/DarthRamone/fireflycrm-bot
ADD . .
RUN go build -o main cmd/app/app.go

FROM alpine
WORKDIR /app

COPY --from=builder /go/src/github.com/DarthRamone/fireflycrm-bot/main .
CMD ["sh", "-c", "./main --token ${TG_TOKEN}"]
