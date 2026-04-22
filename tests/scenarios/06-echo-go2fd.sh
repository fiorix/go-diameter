#!/bin/bash
# 06-echo-go2fd: go echo client sends N Test-Requests to FD's test_app
# running in server mode; expects Result-Code=2001 on each.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# Render test_app.conf (mode=server, codes match tests/dicts/test_app.xml).
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_REALM="${REALM}" \
TEST_APP_MODE=2 TEST_APP_VENDOR=999999 TEST_APP_APP=123456 \
TEST_APP_CMD=234567 TEST_APP_AVP=345678 TEST_APP_LONG_AVP=345679 \
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

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-echo-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-host "${FD_IDENTITY}" -peer-realm "${REALM}"  -dict "${CT_DICT}/test_app.xml"  -n 100 -timeout 15s)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
wait "${go_pid}"

assert_in_log "${go_log}" "done: ok=100 nok=0"
