# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

"Steve" (module `github.com/abolfazlnorzad/steve`) is a Go e-commerce backend: a single Echo HTTP server plus
background schedulers, backed by MongoDB and Redis, with S3-compatible object storage (Arvan/Liara) for media,
Zarinpal for payments, MelliPayamak for SMS, and a Telegram bot for admin/dev notifications. It serves both a
public storefront API and a backoffice (admin panel) API from the same binary.

## Common commands

```bash
go build -o main main.go       # build the server binary (entrypoint: main.go)
go run main.go                 # run locally (reads config.yml + ST_* env vars, needs Mongo/Redis reachable)
go vet ./...                   # static checks
go mod tidy                    # sync go.mod/go.sum
```

There are no automated tests in this repo (`*_test.go` count is zero) and no Makefile/linter config — don't
assume `go test` coverage exists for a change; verify behavior by reading the code path and, if needed, running
the server against local Mongo/Redis (see `docker-compose.yml` for a local Mongo/Redis/app stack, or
`env-tmp.sh` for the shape of required `ST_*` env vars).

Docker: `Dockerfile` is a two-stage build; the final binary needs `CGO_ENABLED=1` because of the `chai2010/webp`
dependency (image encoding), so cross-compiling without CGO will break media handling.

## Configuration

Config is loaded once via `koanf` in `config/loader.go`, in this precedence order: struct defaults
(`config/default.go`) → `config.yml` → environment variables. Env vars must be prefixed `ST_` and use `__` as the
nested-key separator, e.g. `ST_MONGO__HOST` maps to `mongo.host` → `Config.Mongo.Host`. `config.C()` returns the
process-wide singleton; `config.New(opt)` builds an independent instance (used for tests/tools). Add new config by
extending the `Config` struct, its `Default()` entry, and `config.yml`.

## Known pitfalls (recurring failure patterns, not one-off bugs)

- **Silent env-var config misses.** koanf's env loading is purely mechanical string-replacement (`ST_` prefix
  stripped, `__` → `.`) with no schema validation — an env var that doesn't fully spell out every nesting level
  just fails to populate the target struct field, with **no error at startup or anywhere else**. E.g. setting
  `ST_TELEGRAM_BOT=<token>` for a `telegram_bot.admin_api_token_bot` field does nothing; it must be
  `ST_TELEGRAM_BOT__ADMIN_API_TOKEN_BOT`. When adding or reviewing config wiring (`docker-compose.yml`,
  `.gitlab-ci.yml`, `env-tmp.sh`, deployment env), check the *exact* nested path against the `Config` struct in
  `config/`, not just the top-level section name.
- **Copy-pasted `validation.When(...)` gates in `validator/<domain>/update.go`.** Optional-field update
  validators commonly wrap several `validation.Field(...)` rules in `validation.When(req.SomeField != nil, ...)`
  guards, and it's easy to copy one field's block for a new field without updating which field the `When`
  condition checks — the new rule then silently never runs unless the *other* field also happens to be set
  (seen with `TotalUsageLimit`'s validation being gated on `MinimumPurchaseAmount != nil` in
  `validator/discount/update.go`). When adding a validated field to an update DTO, double-check the `When`
  condition tests that same field, and prefer flattening truly-independent optional-field checks instead of
  nesting them under one shared guard.

## Architecture: layered, per-domain vertical slices

The codebase is organized as one directory per architectural layer, and inside each layer one subpackage per
business domain (product, order, cart, discount, brand, category, user, backoffice, payment, sell, repair,
comment, rating, contactus, banner, shipping, article, attribute, acl, search, ...). To understand or extend a
feature, read the same domain name across these layers:

```
entity/          plain domain structs persisted to Mongo (entity.Product, entity.Order, ...)
param/<domain>/  request/response DTOs (CreateXRequest, ListResult, filters) — the shape validators and
                 handlers speak, kept separate from entity/ so storage shape can diverge from API shape
validator/<domain>/  ozzo-validation-based request validation; returns richerror on failure; often depends
                      on other domains' repos/services to check referential integrity (e.g. product validator
                      checks brand/category exist)
service/<domain>/    business logic. Domains commonly expose two service constructors:
                        NewBackofficeService(...)  — admin-facing rules
                        NewStorefrontService(...)  — public-facing rules (e.g. only published/non-deleted items)
                      Services depend on other services/repos via interfaces passed into New*, wired in main.go.
repository/mongo/<domain>/   Mongo persistence, split into one file per operation (add.go, update.go, delete.go,
                              get.go, all.go, ...) sharing a `DB` struct with a `collection` field (see db.go)
repository/redis/<kind>/     Redis-backed repos (e.g. otp)
delivery/http_server/<domain>/  Echo handlers, one file per operation, plus:
                                   handler.go — Handler struct + New(...) wiring deps
                                   routes.go  — SetRoutes(*echo.Echo) registers this domain's routes/middleware
scheduler/<domain>/   background jobs run in goroutines from main.go (order/cart expiry sweeps)
adapter/              outbound integrations: gateway (Zarinpal), mellipayamak (SMS), redis, storage (S3/Liara),
                       telegrambot
```

`main.go` is the composition root: it constructs every repo → service → validator → handler in dependency order
and passes them all into `httpserver.New(...)`. When adding a new domain, follow the existing wiring pattern
there and add a matching block; when adding an operation to an existing domain, only the relevant layer files
need touching (e.g. a new product filter usually only needs `repository/mongo/product/all.go` +
`param/product` + `service/product`, not new files in every layer).

`delivery/http_server/server.go` builds the Echo router, installs global middleware (CORS, request-ID,
structured request logging via `pkg/logger`, panic recovery, and the central error handler), then calls
`SetRoutes` on every handler. Route registration order matters only in that each domain owns its own path
prefix (`/product`, `/order`, `/backoffice/order`, `/storefront/...`, etc.) — storefront-only endpoints
typically live at `<domain>/storefront/...` inside the same handler/routes file as the backoffice endpoints for
that domain, distinguished by which middleware (if any) is attached.

### Auth & ACL

Two independent JWT-based identity kinds share the same `service/auth` package and `pkg/richerror` error model,
distinguished by `claims.Kind`:
- `middleware.Auth` (`delivery/http_server/middleware/auth.go`) — storefront users (`kind == "user"`)
- `backofficemiddleware.BackofficeAuth` (`delivery/http_server/middleware/backoffice`) — admin/backoffice users

Backoffice routes additionally chain `backofficemiddleware.ACL(aclSvc, key.<Permission>)`, checking the
authenticated backoffice user's role permissions (`service/acl`, `entity/acl.go`) against permission constants
defined in `pkg/key`. When adding a protected backoffice endpoint, add/reuse a `key.<Permission>` constant and
wire it into `service/acl`'s seeded permission set (`internal/seeder`) if it's new.

Both auth flows issue access/refresh tokens (constants in `config/constant.go`) and support token blocklisting
via Redis (`CheckTokenIsInBlockList`) for logout/invalidation.

### Errors & responses

- `pkg/richerror`: the only error type handlers/services should construct deliberately. It carries an
  operation name (`op`, conventionally `"<package>.<Func>"`), a `Kind` (Invalid/Forbidden/NotFound/Unexpected/
  BadRequest/Unauthorized — mapped to HTTP status in the central error handler), an i18n `messageKey`, optional
  template data, and optional `meta` (e.g. `{"fields": validation.Errors}` for field-level validation errors).
- `delivery/http_server/middleware/error_handler.go` (`NewCentralErrorHandler`) is the single place that turns
  any error into an `APIResponse` (`delivery/http_server/response`) and an HTTP status — it unwraps
  `richerror`, translates the message key via `pkg/i18n` using the `Accept-Language` header (`en`/`fa`, default
  `fa`), and translates per-field validation errors (including translating the field *name* via an `attr.<field>`
  i18n key). Non-richerror / non-5xx-safe errors collapse to a generic `global.something_wrong` message so
  internals never leak to clients.
- i18n message catalogs live under `config/locales`; add new user-facing message keys there in both languages,
  not just in code.

### Conventions to follow when extending a domain

- New Mongo repo operations get their own file under `repository/mongo/<domain>/` (e.g. `get.go`, `update.go`)
  rather than being added to an existing file — mirror the existing per-operation file split.
- New HTTP endpoints likewise get their own file under `delivery/http_server/<domain>/`, registered in that
  domain's `routes.go`.
- Prefer interfaces at service boundaries (services take other services'/repos' interfaces, not concrete
  structs) so wiring stays swappable from `main.go`.
- IDs are ULIDs (`github.com/oklog/ulid/v2`), not Mongo ObjectIDs, for externally-facing entity IDs.
