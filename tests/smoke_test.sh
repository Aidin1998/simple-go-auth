#!/usr/bin/env bash
set -euo pipefail

HOST=${SERVICE_HOST:-localhost:443}

echo "Waiting for service..."
for i in {1..12}; do
  if curl --insecure -sf "https://$HOST/health"; then
    echo "✅ /health OK"
    break
  fi
  sleep 5
done

echo "Testing signup flow..."
curl --insecure -sf -XPOST "https://$HOST/signup" \
  -H 'Content-Type: application/json' \
  -d '{"username":"smoketest","password":"Test1234","email":"test@x.com"}'

echo "Smoke tests passed!"
