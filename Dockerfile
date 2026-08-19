# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirror.arvancloud.ir/g' /etc/apk/repositories

RUN apk add --no-cache git build-base

# تنظیم دایرکتوری کاری
WORKDIR /app

# کپی فایل‌های وابسته و دانلود آن‌ها
COPY go.mod go.sum ./
# proxy.golang.org و sum.golang.org گوگل هاست می‌کنن و از سرور ایرانی وصل شدن
# بهشون intermittent / گاهی 403 هست (دیده شد حین بیلد واقعی). goproxy.cn و
# goproxy.io هردو تست شدن و قابل دسترسن، پس اول اونا رو امتحان می‌کنه.
# GOSUMDB هم خاموشه چون go.sum از قبل تو ریپو هست و همون verify می‌کنه؛ لازم
# نیست sum.golang.org (که اونم گوگل هست) رو صدا بزنه.
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GOSUMDB=off
RUN go mod download

# کپی تمام کدهای پروژه
COPY . .

# ✅ اصلاح کلیدی: فعال‌سازی CGO برای پکیج webp
# پکیج chai2010/webp برای کامپایل بخش‌های C خود به CGO نیاز دارد
RUN CGO_ENABLED=1 GOOS=linux go build -o main main.go

# Stage 2: Final lightweight image
FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirror.arvancloud.ir/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates tzdata libc6-compat
ENV TZ=Asia/Tehran

WORKDIR /root/

# کپی فایل باینری از مرحله قبل
COPY --from=builder /app/main .

# کپی فایل‌های ترجمه (i18n) - این‌ها بخشی از کد هستند و در گیت هست
COPY --from=builder /app/config ./config

# توجه: config.yml عمداً اینجا کپی نمی‌شود. این فایل در .gitignore است
# (حاوی سکرت‌های واقعی است) و در چک‌اوت تازه‌ی گیت اصلاً وجود ندارد؛
# اگر اینجا COPY می‌شد، بیلد داکر همیشه fail می‌کرد. config.yml باید
# در زمان اجرا با یک bind mount در docker-compose.yml فراهم شود
# (نگاه کن به سرویس app در docker-compose.yml).

# باز کردن پورت برنامه
EXPOSE 7777

# اجرای برنامه
CMD ["./main"]