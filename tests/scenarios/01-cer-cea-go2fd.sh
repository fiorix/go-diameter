#!/bin/bash
# 01-cer-cea-go2fd: go-diameter client connects to freeDiameterd,
# completes the CER/CEA handshake, sends DPR, exits cleanly.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"

: "${SCEN_DIR:?SCEN_DIR must be set by run.sh}"

ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# Render FD config (FD is the server/responder).
FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
    envsubst < "${TESTS_DIR}/configs/fd-base.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"

fd_log="${SCEN_DIR}/fd.log"
go_log="${SCEN_DIR}/go-client.log"

# Start FD in background.
wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
trap 'kill_pids ${fd_pid:-} ${go_pid:-}' EXIT

wait_for_line "${fd_log}" "freeDiameterd daemon initialized" 15 \
    || { log "FD didn't start"; exit 1; }

# Run go client.
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-host "${FD_IDENTITY}" -peer-realm "${REALM}"  -timeout 5s -dpr)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
wait "${go_pid}"

# Verify.
assert_in_log "${go_log}" "CER/CEA handshake OK"
assert_in_log "${go_log}" "DPR sent"
assert_in_log "${fd_log}" "CER|Capabilities-Exchange-Request"
assert_in_log "${fd_log}" "CEA|Capabilities-Exchange-Answer"
