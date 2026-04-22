#!/bin/bash
# 05-cer-cea-reject: go client advertises a mismatched Origin-Realm
# against FD's acl_wl; FD should reject with DIAMETER_UNKNOWN_PEER or
# simply drop the connection.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# FD is configured with only fdiam.test.local as a peer; go will claim
# a different Origin-Host that is NOT in acl_wl.
FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
    envsubst < "${TESTS_DIR}/configs/fd-base.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"

fd_log="${SCEN_DIR}/fd.log"
go_log="${SCEN_DIR}/go-client.log"

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
trap 'kill_pids ${fd_pid:-} ${go_pid:-}' EXIT
wait_for_line "${fd_log}" "freeDiameterd daemon initialized" 15

# Use a bogus Origin-Host that isn't in FD's ConnectPeer or acl_wl list.
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "rogue.evil.example" -origin-realm "evil.example"  -timeout 3s -dpr=false)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!

# sdme exec does not propagate child exit codes reliably, so assert
# on log content. FD should return 3010 DIAMETER_UNKNOWN_PEER in its
# CEA and the go client logs the failure.
wait "${go_pid}" 2>/dev/null || true

grep -qE "3010|DIAMETER_UNKNOWN_PEER|failed Result-Code" \
    "${fd_log}" "${go_log}" \
    || { log "expected 3010/UNKNOWN_PEER indication"; \
         tail -n 40 "${fd_log}" >&2; tail -n 40 "${go_log}" >&2; exit 1; }
