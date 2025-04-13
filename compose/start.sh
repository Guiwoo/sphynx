#! /bin/bash

set -e

echo "[+] Stopping current containers..."
docker compose -f infra.yml down

echo "[+] Starting updated containers..."
docker compose -f infra.yml up -d --build

echo "[+] Done. ✅"