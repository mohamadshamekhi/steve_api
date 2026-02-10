# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/uk.alpinelinux.org/g' /etc/apk/repositories

RUN apk add --no-cache git build-base

# تنظیم دایرکتوری کاری
WORKDIR /app

# کپی فایل‌های وابسته و دانلود آن‌ها
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# کپی تمام کدهای پروژه
COPY . .

# ✅ اصلاح کلیدی: فعال‌سازی CGO برای پکیج webp
# پکیج chai2010/webp برای کامپایل بخش‌های C خود به CGO نیاز دارد
RUN CGO_ENABLED=1 GOOS=linux go build -o main main.go

# Stage 2: Final lightweight image
FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/uk.alpinelinux.org/g' /etc/apk/repositories
# ✅ نصب ca-certificates برای ارتباطات امن و libc6-compat برای اجرای باینری‌های CGO در Alpine
RUN apk --no-cache add ca-certificates tzdata libc6-compat
ENV TZ=Asia/Tehran

WORKDIR /root/

# کپی فایل باینری از مرحله قبل
COPY --from=builder /app/main .

# کپی فایل تنظیمات (yml)
COPY --from=builder /app/config.yml .
COPY --from=builder /app/config ./config

# باز کردن پورت برنامه
EXPOSE 7777

# اجرای برنامه
CMD ["./main"]