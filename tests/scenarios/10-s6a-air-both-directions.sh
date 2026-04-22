#!/bin/bash
# 10-s6a-air-both-directions: go client sends AIR to FD (test_app cmd=318
# in server mode); then FD test_app client pings go s6a-server with AIR.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_REALM="${REALM}" \
TEST_APP_MODE=3 TEST_APP_VENDOR=10415 TEST_APP_APP=16777251 \
TEST_APP_CMD=318 TEST_APP_AVP=345678 TEST_APP_LONG_AVP=345679 \
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
goc_log="${SCEN_DIR}/go-client.log"

# 1) Start go s6a-server (for FD->go direction)
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-s6a-server"  -addr ":${GO_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}")
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
gos_pid=$!
trap 'kill_pids ${gos_pid:-} ${fd_pid:-} ${goc_pid:-}; stop_processes_in "${GO_CT}" "diam-s6a"' EXIT
wait_for_line "${go_log}" "s6a server listening" 10

# 2) Start FD (both peer roles)
wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
wait_for_line "${fd_log}" "STATE_OPEN" 15

# 3) Fire FD test_app client 3x at go server
for _ in 1 2 3; do
    sdme exec "${FD_CT}" -- /usr/bin/pkill -USR1 -x freeDiameterd || true
    sleep 0.3
done

# 4) go client sends AIR to FD
sdme exec "${GO_CT}" -- "${CT_BIN}/diam-s6a-client" \
    -addr "127.0.0.1:${FD_PORT}" \
    -origin-host "go2.test.local" -origin-realm "${REALM}" \
    -peer-host "${FD_IDENTITY}" -peer-realm "${REALM}" \
    -timeout 10s \
    >"${goc_log}" 2>&1 || true

# 5) Verify both directions saw 2001s.
fd_2001=$(grep -cE "RESULT_CODE.*2001|Result-Code.*2001" "${fd_log}" || true)
[[ "${fd_2001}" -ge 3 ]] || { log "FD side: expected >=3 2001, got ${fd_2001}"; exit 1; }
grep -qE "AIA Result-Code=2001|ULA Result-Code=2001" "${goc_log}" \
    || { log "go client didn't see 2001"; tail -n 40 "${goc_log}" >&2; exit 1; }
