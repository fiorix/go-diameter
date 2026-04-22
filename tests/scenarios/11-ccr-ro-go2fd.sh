#!/bin/bash
# 11-ccr-ro-go2fd: go CC client sends CCR-I/U/T (Ro/Gy, App-Id 4) to
# FD configured to echo CC command (272) with Result-Code=2001.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_REALM="${REALM}" \
TEST_APP_MODE=2 TEST_APP_VENDOR=999999 TEST_APP_APP=4 \
TEST_APP_CMD=272 TEST_APP_AVP=345678 TEST_APP_LONG_AVP=345679 \
    envsubst < "${TESTS_DIR}/configs/test_app.conf.tmpl" \
    > "${SCEN_DIR}/test_app.conf"

FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
FD_TESTAPP_CONF="${ct_scen}/test_app.conf" \
    envsubst < "${TESTS_DIR}/configs/fd-echo.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"

fd_log="${SCEN_DIR}/fd.log"
go_log="${SCEN_DIR}/go-client.log"

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
trap 'kill_pids ${fd_pid:-} ${go_pid:-}' EXIT
wait_for_line "${fd_log}" "freeDiameterd daemon initialized" 15

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-cc-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-realm "${REALM}"  -app-id 4 -timeout 10s)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
wait "${go_pid}"

assert_in_log "${go_log}" "OK"
count=$(grep -cE "CCA Result-Code=2001" "${go_log}" || true)
[[ "${count}" -eq 3 ]] || { log "expected 3 CCA 2001, got ${count}"; exit 1; }
