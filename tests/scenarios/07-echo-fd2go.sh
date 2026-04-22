#!/bin/bash
# 07-echo-fd2go: FD test_app (mode=client) drives go echo server. SIGUSR1
# is sent to freeDiameterd to trigger a test request; we repeat N times.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"
N=20

FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_REALM="${REALM}" \
TEST_APP_MODE=1 TEST_APP_VENDOR=999999 TEST_APP_APP=123456 \
TEST_APP_CMD=234567 TEST_APP_AVP=345678 TEST_APP_LONG_AVP=345679 \
    envsubst < "${TESTS_DIR}/configs/test_app.conf.tmpl" \
    > "${SCEN_DIR}/test_app.conf"

FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
FD_TESTAPP_CONF="${ct_scen}/test_app.conf" \
    envsubst < "${TESTS_DIR}/configs/fd-echo.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"

go_log="${SCEN_DIR}/go-server.log"
fd_log="${SCEN_DIR}/fd.log"

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-echo-server"  -addr ":${GO_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -dict "${CT_DICT}/test_app.xml")
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
trap 'kill_pids ${go_pid:-} ${fd_pid:-}; stop_processes_in "${GO_CT}" diam-echo-server' EXIT
wait_for_line "${go_log}" "echo server listening" 10

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
wait_for_line "${fd_log}" "STATE_OPEN" 15

# test_app client fires one request per SIGUSR1.
for _ in $(seq 1 ${N}); do
    sdme exec "${FD_CT}" -- /usr/bin/pkill -USR1 -x freeDiameterd || true
    sleep 0.2
done
sleep 2

# Verify FD recorded N successful answers.
ok=$(grep -cE "RESULT_CODE.*2001|Result-Code.*2001" "${fd_log}" || true)
[[ "${ok}" -ge ${N} ]] || { log "expected >=${N} 2001 answers, got ${ok}"; exit 1; }
