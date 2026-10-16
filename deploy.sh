#!/usr/bin/env bash
#
# deploy.sh - ship the steve API to the server and restart it.
#
# Builds ON the server itself (not locally) on purpose: this repo is developed
# on Apple Silicon Macs (arm64), the server is a standard amd64 Linux VPS, and
# CGO_ENABLED=1 (needed for chai2010/webp) makes cross-compiling from a local
# build genuinely painful. rsync + remote `docker build` sidesteps all of
# that, and matches the pattern already used for the storefront/backoffice
# services on this same server.
#
# What it does, in order:
#   1. sanity checks (sshpass/rsync/ssh present, credentials set, server reachable)
#   2. `go build main.go` locally - fails fast before touching the network at all
#   3. rsync the repo to the server (excluding config.yml/.env/.git/... - see
#      RSYNC_EXCLUDES below; those hold secrets or are server-local state and
#      are never overwritten by a deploy)
#   4. confirms config.yml and .env already exist on the server (first deploy
#      needs them placed there once, see "first-time server setup" below)
#   5. `mongodump`s steve_db to REMOTE_ROOT/backups/ and verifies the archive
#      is non-empty and passes `gzip -t` before continuing (skipped if this is
#      the first deploy on this server, i.e. no steve_db yet to back up)
#   6. `docker build` on the server, then `docker compose up -d --force-recreate`
#   7. a smoke test against the app port, then prints recent logs on failure
#
# Usage:
#   ./deploy.sh                  # deploy using deploy.env (or env vars) for creds
#   ./deploy.sh --skip-build-check # skip the local `go build` compile check
#   ./deploy.sh --push-config     # also (re-)upload local config.yml to the server
#
# Credentials: copy deploy.env.example to deploy.env (gitignored) and fill it
# in, or export SERVER_USER/SERVER_IP/SERVER_PASS/SERVER_PORT yourself.
# Leave SERVER_PASS unset to use a plain SSH key/agent instead of sshpass.
#
# First-time server setup (only needed once per server):
#   1. docker + the docker compose plugin installed on the server
#   2. `scp config.yml root@<ip>:/root/project/api/config.yml` (or --push-config)
#   3. `scp .env.example root@<ip>:/root/project/api/.env` then fill in real values
#
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

# ---------------------------------------------------------------------------
# configuration
# ---------------------------------------------------------------------------
[[ -f deploy.env ]] && source deploy.env

SERVER_USER="${SERVER_USER:-root}"
SERVER_IP="${SERVER_IP:-}"
SERVER_PASS="${SERVER_PASS:-}"
SERVER_PORT="${SERVER_PORT:-22}"
REMOTE_ROOT="${REMOTE_ROOT:-/root/project/api}"
IMAGE_NAME="${IMAGE_NAME:-steve-api:latest}"
SKIP_BUILD_CHECK=0
PUSH_CONFIG=0

usage() { sed -n '2,29p' "$0" | sed 's/^# \{0,1\}//'; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-build-check) SKIP_BUILD_CHECK=1; shift ;;
    --push-config) PUSH_CONFIG=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage; exit 1 ;;
  esac
done

[[ -n "$SERVER_IP" ]] || { echo "SERVER_IP is not set - copy deploy.env.example to deploy.env and fill it in" >&2; exit 1; }

if [[ -n "$SERVER_PASS" ]]; then
  command -v sshpass >/dev/null || { echo "sshpass not installed (brew install sshpass), or unset SERVER_PASS to use SSH keys instead" >&2; exit 1; }
  SSH_PREFIX="sshpass -p ${SERVER_PASS} ssh"
  RSYNC_SSH="sshpass -p ${SERVER_PASS} ssh"
else
  SSH_PREFIX="ssh"
  RSYNC_SSH="ssh"
fi
SSH="$SSH_PREFIX -p ${SERVER_PORT} -o StrictHostKeyChecking=no -o ConnectTimeout=10 -o ServerAliveInterval=15 -o ServerAliveCountMax=8 ${SERVER_USER}@${SERVER_IP}"
REMOTE="${SERVER_USER}@${SERVER_IP}"

RSYNC_EXCLUDES=(
  --exclude '.git' --exclude '.idea' --exclude '.DS_Store'
  --exclude '.bin' --exclude 'logs' --exclude 'main'
  --exclude 'config.yml' --exclude '.env' --exclude 'deploy.env'
  --exclude 'env-tmp.sh'
)

# ---------------------------------------------------------------------------
# logging helpers
# ---------------------------------------------------------------------------
step() { printf '\n\033[1;34m==>\033[0m \033[1m%s\033[0m\n' "$1"; }
ok()   { printf '\033[1;32m✓\033[0m %s\n' "$1"; }
fail() { printf '\033[1;31m✗ %s\033[0m\n' "$1" >&2; exit 1; }

trap 'echo; fail "deploy failed (see above) - server was left as-is by the step that failed"' ERR

# ---------------------------------------------------------------------------
# 1. preflight
# ---------------------------------------------------------------------------
step "Preflight checks"
command -v rsync >/dev/null || fail "rsync is not installed locally"
[[ -f Dockerfile ]] || fail "Dockerfile not found - run this from the repo root"
[[ -f docker-compose.yml ]] || fail "docker-compose.yml not found - run this from the repo root"

$SSH exit 2>/dev/null || fail "cannot reach ${REMOTE}:${SERVER_PORT} over SSH - check deploy.env"
ok "SSH to ${REMOTE} works"

$SSH "command -v docker >/dev/null && docker compose version >/dev/null" \
  || fail "docker (with the compose plugin) is not installed on ${REMOTE}"
ok "docker + compose plugin present on server"

# ---------------------------------------------------------------------------
# 2. compile check (fail fast before touching the network)
# ---------------------------------------------------------------------------
# Deliberately `go build main.go`, not `go vet ./...`: vet type-checks every
# package in the module, including ones main.go doesn't even import (dead
# code can and does bit-rot in this repo - see service/variant), and flags
# style-only issues (unkeyed struct literals, slog arg pairing) that don't
# stop a build. Building main.go is exactly what the Dockerfile does, so it's
# the accurate "will the image build" signal without unrelated noise.
if [[ $SKIP_BUILD_CHECK -eq 0 ]]; then
  step "go build (compile check)"
  if command -v go >/dev/null; then
    go build -o /dev/null main.go || fail "go build failed - fix the errors above (or re-run with --skip-build-check)"
    ok "main.go and its dependencies compile"
  else
    echo "  go toolchain not found locally, skipping (remote docker build will still catch compile errors)"
  fi
fi

# ---------------------------------------------------------------------------
# 3. rsync the repo to the server
# ---------------------------------------------------------------------------
step "Syncing code to ${REMOTE}:${REMOTE_ROOT}"
$SSH "mkdir -p '${REMOTE_ROOT}'"
rsync -avz --delete --checksum "${RSYNC_EXCLUDES[@]}" -e "$RSYNC_SSH -p ${SERVER_PORT} -o StrictHostKeyChecking=no -o ServerAliveInterval=15 -o ServerAliveCountMax=8" \
  ./ "${REMOTE}:${REMOTE_ROOT}/"
ok "code synced"

# ---------------------------------------------------------------------------
# 4. make sure server-side secrets already exist (config.yml / .env)
# ---------------------------------------------------------------------------
step "Checking server-side config"

if [[ $PUSH_CONFIG -eq 1 ]]; then
  [[ -f config.yml ]] || fail "--push-config given but no local config.yml found"
  echo "  uploading local config.yml to ${REMOTE}:${REMOTE_ROOT}/config.yml"
  echo "  (this puts real secrets on the wire - make sure that's intended)"
  scp_prefix=(); [[ -n "$SERVER_PASS" ]] && scp_prefix=(sshpass -p "$SERVER_PASS")
  "${scp_prefix[@]}" scp -P "${SERVER_PORT}" -o StrictHostKeyChecking=no config.yml "${REMOTE}:${REMOTE_ROOT}/config.yml"
  ok "config.yml uploaded"
fi

$SSH "[ -f '${REMOTE_ROOT}/config.yml' ]" || fail "no config.yml on the server at ${REMOTE_ROOT}/config.yml.
This file holds real secrets and is intentionally never auto-generated. Either place it
there yourself once, or re-run this script with --push-config to upload your local one."
ok "config.yml present on server"

$SSH "[ -f '${REMOTE_ROOT}/.env' ]" || fail "no .env on the server at ${REMOTE_ROOT}/.env.
Copy .env.example there once and fill in real credentials."
ok ".env present on server"

# sanity check: the mongo data volume must already exist so --force-recreate
# below reuses it instead of silently starting from an empty database.
MONGO_EXISTS=0
if $SSH "docker ps -a --format '{{.Names}}' | grep -qx steve_db"; then
  MONGO_EXISTS=1
  if ! $SSH "docker volume ls --format '{{.Name}}' | grep -qx api_mongo_data"; then
    fail "steve_db container exists but the expected 'api_mongo_data' volume does not.
Refusing to continue - recreating steve_db right now could start from an empty database.
Investigate manually (docker inspect steve_db) before re-running."
  fi
  ok "existing api_mongo_data volume found - will be reused"
fi

# ---------------------------------------------------------------------------
# 5. back up Mongo before touching anything. --force-recreate below restarts
# steve_db too (it reuses the same volume, but this is the safety net in case
# anything about that goes wrong).
# ---------------------------------------------------------------------------
if [[ $MONGO_EXISTS -eq 1 ]]; then
  step "Backing up Mongo data on ${REMOTE}"
  BACKUP_DIR="${REMOTE_ROOT}/backups"
  BACKUP_PATH="${BACKUP_DIR}/mongo_$(date +%Y%m%d_%H%M%S).archive.gz"

  $SSH "mkdir -p '${BACKUP_DIR}'"
  $SSH "docker exec steve_db sh -c 'mongodump --username \"\$MONGO_INITDB_ROOT_USERNAME\" --password \"\$MONGO_INITDB_ROOT_PASSWORD\" --authenticationDatabase admin --archive --gzip' > '${BACKUP_PATH}'" \
    || fail "mongodump failed - aborting before touching the running database"

  # a backup that's zero bytes or fails gzip's integrity check is worse than no
  # backup at all (false confidence) - refuse to continue on either.
  BACKUP_SIZE=$($SSH "stat -c%s '${BACKUP_PATH}' 2>/dev/null")
  [[ "${BACKUP_SIZE:-0}" -gt 0 ]] || fail "backup file ${BACKUP_PATH} is empty - aborting, nothing was restarted"
  $SSH "gzip -t '${BACKUP_PATH}'" || fail "backup file ${BACKUP_PATH} failed gzip integrity check - aborting, nothing was restarted"
  ok "backup verified OK: ${BACKUP_PATH} (${BACKUP_SIZE} bytes)"

  # keep the last 10 backups only, so this doesn't grow the server's disk forever
  $SSH "cd '${BACKUP_DIR}' && ls -t mongo_*.archive.gz 2>/dev/null | tail -n +11 | xargs -r rm --" || true
else
  step "Backing up Mongo data"
  echo "  no existing steve_db container - nothing to back up yet (first deploy on this server)"
fi

# ---------------------------------------------------------------------------
# 6. build + restart on the server
# ---------------------------------------------------------------------------
step "Building image on ${REMOTE}"
$SSH "cd '${REMOTE_ROOT}' && docker build -t '${IMAGE_NAME}' ."
ok "image built on server"

step "Restarting the stack on ${REMOTE}"
$SSH "cd '${REMOTE_ROOT}' && docker compose up -d --force-recreate"
ok "docker compose up -d --force-recreate done"

# ---------------------------------------------------------------------------
# 7. smoke test
# ---------------------------------------------------------------------------
step "Smoke test"
# Poll and fetch logs inside ONE ssh session: startup (Mongo/Redis connect,
# index builds) can take longer than a few seconds, and opening a second
# connection right after a failure has been rejected by the server's auth
# throttling, which hid the logs we need.
if $SSH "for i in \$(seq 1 30); do
    curl -fsS -o /dev/null -m 5 http://localhost:7777/; ec=\$?
    [ \$ec -eq 0 -o \$ec -eq 22 ] && exit 0
    sleep 2
  done
  echo '  app did not answer on port 7777 within 60s - container state & recent logs:'
  cd '${REMOTE_ROOT}' && docker compose ps app && docker compose logs --tail=150 app
  exit 1"; then
  ok "app is answering on port 7777"
else
  fail "smoke test failed"
fi

step "Pruning dangling images on server"
$SSH "docker image prune -f" >/dev/null || true

step "Deploy finished: ${IMAGE_NAME} is live on ${REMOTE}"
