#!/bin/sh
set -eu

cat > /usr/share/nginx/html/config.js <<EOF
window.__OTASIGN_CONFIG__ = {
  apiBaseUrl: "${OTASIGN_API_BASE_URL:-http://localhost:8080}"
};
EOF
