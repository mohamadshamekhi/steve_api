#!/usr/bin/env bash
#
# server-deploy.sh - runs ON the server (called by GitHub Actions via SSH).
# Usage: ./server-deploy.sh /path/to/repo
#
# One-time server setup before this can run:
#   1. git clone git@github.com:mohamadshamekhi/steve_api.git <DEPLOY_PATH>
#   2. Place config.yml and .env in <DEPLOY_PATH> (never committed)
#   3. docker + compose plugin installed
#
set -euo pipefail

DEPLOY_PATH="${1:-$(cd "$(dirname "$0")/.." && pwd)}"
IMAGE_NAME="${IMAGE_NAME:-steve-api:latest}"

step() { printf '\n\033[1;34m==>\033[0m \033[1m%s\033[0m\n' "$1"; }
ok()   { printf '\033[1;32m✓\033[0m %s\n' "$1"; }
fail() { printf '\033[1;31m✗ %s\033[0m\n' "$1" >&2; exit 1; }

trap 'fail "deploy failed (see above)"' ERR

cd "$DEPLOY_PATH"

[[ -f Dockerfile ]]         || fail "Dockerfile not found at $DEPLOY_PATH"
[[ -f docker-compose.yml ]] || fail "docker-compose.yml not found at $DEPLOY_PATH"
[[ -f config.yml ]]         || fail "config.yml missing - place it manually on the server once"
[[ -f .env ]]               || fail ".env missing - copy .env.example and fill in secrets"

# ---------------------------------------------------------------------------
# 1. pull latest code
# ---------------------------------------------------------------------------
step "git pull"
git fetch --prune origin
git reset --hard origin/main
ok "code updated to $(git rev-parse --short HEAD)"

# ---------------------------------------------------------------------------
# 2. mongo backup (skip if steve_db container not yet running)
# ---------------------------------------------------------------------------
if docker ps -a --format '{{.Names}}' | grep -qx steve_db; then
  step "Mongo backup"
  BACKUP_DIR="$DEPLOY_PATH/backups"
  BACKUP_PATH="$BACKUP_DIR/mongo_$(date +%Y%m%d_%H%M%S).archive.gz"
  mkdir -p "$BACKUP_DIR"
  docker exec steve_db sh -c \
    'mongodump --username "$MONGO_INITDB_ROOT_USERNAME" --password "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin --archive --gzip' \
    > "$BACKUP_PATH" || fail "mongodump failed"
  BACKUP_SIZE=$(stat -c%s "$BACKUP_PATH" 2>/dev/null || echo 0)
  [[ "${BACKUP_SIZE:-0}" -gt 0 ]] || fail "backup is empty"
  gzip -t "$BACKUP_PATH" || fail "backup failed integrity check"
  ok "backup OK: $BACKUP_PATH ($BACKUP_SIZE bytes)"
  # keep last 10
  cd "$BACKUP_DIR" && ls -t mongo_*.archive.gz 2>/dev/null | tail -n +11 | xargs -r rm -- || true
  cd "$DEPLOY_PATH"
else
  step "No existing steve_db - skipping backup (first deploy)"
fi

# ---------------------------------------------------------------------------
# 3. build + restart
# ---------------------------------------------------------------------------
step "docker build"
docker build -t "$IMAGE_NAME" .
ok "image built"

step "docker compose up"
docker compose up -d --force-recreate
ok "containers restarted"

# ---------------------------------------------------------------------------
# 4. smoke test
# ---------------------------------------------------------------------------
step "Smoke test (port 7777)"
for i in $(seq 1 60); do
  curl -fsS -o /dev/null -m 5 http://localhost:7777/health && { ok "app is up"; break; } || true
  [[ $i -eq 60 ]] && {
    echo "app did not answer within 120s:"
    docker compose logs --tail=80 app
    fail "smoke test failed"
  }
  sleep 2
done

step "Pruning dangling images"
docker image prune -f >/dev/null || true

step "Deploy done: $IMAGE_NAME is live"
