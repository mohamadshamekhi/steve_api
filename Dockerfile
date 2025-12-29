FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go


FROM alpine:latest


RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Tehran

WORKDIR /root/


COPY --from=builder /app/main .


COPY --from=builder /app/config.yml .


EXPOSE 7777


CMD ["./main"]