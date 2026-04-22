#!/bin/bash
# Generate a throwaway self-signed CA + server cert for FreeDiameter.
# FD requires TLS_Cred be configured even with No_TLS, but never reads
# the cert CN because TLS is never negotiated. Run once; commit the
# output or let the bash harness regenerate on demand.

set -euo pipefail

OUT_DIR="${1:-$(dirname "$0")}"
CN="${2:-fdiam.test.local}"

cd "${OUT_DIR}"

# FD v1.5.0 validates the cert CN matches its Identity at startup even
# when peers are set to No_TLS. Include both test FQDNs as SAN so the
# same cert works when either side is FD.
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem \
    -days 3650 -nodes -subj "/CN=${CN}" \
    -addext "subjectAltName=DNS:fdiam.test.local,DNS:godiam.test.local"
cp cert.pem ca.pem

chmod 0644 cert.pem ca.pem
chmod 0600 key.pem

echo "wrote ${OUT_DIR}/{cert,key,ca}.pem"
