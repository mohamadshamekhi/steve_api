# Steve API — Postman Collection

مستندسازی کامل API بک‌اند Steve (فروشگاه + پنل بک‌آفیس) به‌صورت یک Postman Collection، تولیدشده با تحلیل استاتیک تمام
فایل‌های `delivery/http_server/**/routes.go` در ریپازیتوری (بدون اجرای سرور و بدون هیچ فراخوانی شبکه‌ای).

این پوشه شامل دو فایل است:

- `steve-api.postman_collection.json` — کالکشن کامل (133 request، در 23 پوشهٔ دامنه‌ای)
- `steve-api.postman_environment.json` — Environment همراه با متغیرهای لازم (`base_url`, `user_token`, `backoffice_token`, و چند شناسهٔ نمونه)

## Import (وارد کردن فایل‌ها)

1. در Postman: **File → Import** و هر دو فایل JSON بالا را انتخاب کنید (یا drag & drop).
2. بعد از Import، از منوی بالا-راست Postman (Environment selector) گزینهٔ **"Steve API"** را انتخاب کنید تا
   متغیرهای محیطی فعال شوند.

## راه‌اندازی سرور محلی

`base_url` به‌صورت پیش‌فرض روی `http://localhost:7777` تنظیم شده که دقیقاً برابر با مقدار `http_server.port` در
`config.yml` است. سرور را طبق دستورالعمل `CLAUDE.md` (`go run main.go`، با Mongo/Redis در دسترس) اجرا کنید.

## گرفتن توکن‌ها (Login)

دو نوع JWT کاملاً مستقل در این API وجود دارد:

### 1. کاربر فروشگاه (Storefront User) — `user_token`

در پوشهٔ **User → Public**:

1. `Send OTP (User Login)` — شمارهٔ موبایل را ارسال کنید تا کد OTP پیامک شود.
2. `Verify OTP (User Login)` — کد دریافتی را همراه شماره موبایل ارسال کنید؛ پاسخ شامل فیلد `token` است.
3. مقدار `token` را کپی کرده و در Environment، در متغیر `user_token` قرار دهید.

هر پوشه‌ای که auth آن روی `{{user_token}}` تنظیم شده (مثلاً Cart → Storefront، Order → Storefront، Comment →
Storefront، Rating، Repair/Sell → Storefront، Payment، بخش‌هایی از User) از این توکن به‌صورت Bearer استفاده
می‌کند.

### 2. کاربر پنل ادمین (Backoffice) — `backoffice_token`

در پوشهٔ **Backoffice Auth → Public**:

1. `Send OTP (Backoffice Login)`
2. `Verify OTP (Backoffice Login)` — پاسخ شامل فیلد `token` است.
3. مقدار `token` را در متغیر `backoffice_token` قرار دهید.

هر پوشهٔ `Backoffice` در کالکشن از این توکن استفاده می‌کند. علاوه بر توکن معتبر، اکثر route های بک‌آفیس یک
دسترسی ACL مشخص (`key.<Permission>` در `pkg/key`) نیز نیاز دارند که در قسمت description هر request قید شده
است؛ کاربر بک‌آفیس شما باید نقش/دسترسی متناظر را داشته باشد وگرنه پاسخ 403 دریافت می‌کنید.

## نکات ساختاری

- برای دامنه‌هایی که هم مخاطب فروشگاه و هم بک‌آفیس دارند، دو زیرپوشهٔ `Storefront` و `Backoffice` (و گاهی
  `Public` برای route های کاملاً بدون auth) ساخته شده؛ auth روی سطح پوشه تنظیم شده تا در تک‌تک request ها
  تکرار نشود (Postman آن را ارث‌بری می‌کند).
- برای PATCH/PUT request ها که فیلدهای اختیاری (pointer) دارند، فقط چند فیلد نمونه در body قرار داده شده — نیازی
  نیست همهٔ فیلدها را پر کنید.
- برخی مقادیر id در body/URL نمونه‌های ثابت (ULID ساختگی مثل `01ARZ3NDEKTSV4RRFFQ69G5FAV`) هستند؛ آن‌ها را با
  idهای واقعی محیط خودتان جایگزین کنید (یا در Environment متغیرهایی مثل `product_id`, `order_id` و... را پر
  کرده و در URL/body دستی جایگزین کنید).
- برای دو مسیر آپلود فایل (`Media → Backoffice/Storefront → Upload File`) بدنهٔ درخواست از نوع
  `multipart/form-data` است؛ فیلد `file` را در Postman با یک فایل واقعی پر کنید.

## نکات مهم دربارهٔ کد که در کالکشن هم مستند شده‌اند

- `GET /contact-us/backoffice` در سورس کد **هیچ middleware احراز هویتی ندارد** (برخلاف الگوی سایر route های
  بک‌آفیس) — این احتمالاً یک oversight است، ولی در کالکشن دقیقاً همان‌طور که در کد هست (بدون auth) مستند شده.
- چند دامنه (`attribute`, `attribute_definition`, `shipping`, `repair`, `sell`) از کلیدهای ACL دامنهٔ دیگر
  استفاده می‌کنند (مثلاً `key.CreateBrand` برای ساخت attribute/shipping، یا `key.UpdateProduct` برای
  repair/sell) — در description هر request قید شده.
- `AddVariant`/`UpdateVariant` محصول (`param/product/backoffice_product_variant.go`) تگ `json` ندارند، پس
  encoding/json جاوااسکریپت-مانند نه، بلکه نام فیلد Go (با حرف بزرگ، مثل `Price`, `Warranty`) را به‌عنوان کلید
  JSON استفاده می‌کند — نمونهٔ body در کالکشن این رفتار واقعی را منعکس می‌کند.

هر دو فایل JSON با `python3 -m json.tool` اعتبارسنجی شده‌اند و تعداد 133 request در کالکشن، دقیقاً با شمارش خط
به خط ثبت route در 22 فایل `routes.go` مطابقت دارد.
