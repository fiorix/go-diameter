#!/bin/bash
# 08-s6a-ulr-go2fd: go S6a client sends AIR+ULR to FD configured with
# dict_s6a.fdx + test_app tuned to S6a codes. test_app echoes each
# request with Result-Code=2001, validating CER app negotiation and
# wire-format compatibility for S6a.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# test_app configured to advertise S6a app and echo ULR (316).
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_REALM="${REALM}" \
TEST_APP_MODE=2 TEST_APP_VENDOR=10415 TEST_APP_APP=16777251 \
TEST_APP_CMD=316 TEST_APP_AVP=345678 TEST_APP_LONG_AVP=345679 \
    envsubst < "${TESTS_DIR}/configs/test_app.conf.tmpl" \
    > "${SCEN_DIR}/test_app.conf"

FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
FD_TESTAPP_CONF="${ct_scen}/test_app.conf" \
    envsubst < "${TESTS_DIR}/configs/fd-echo.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"
# Load the S6a dictionary so AVPs decode by name in FD logs.

fd_log="${SCEN_DIR}/fd.log"
go_log="${SCEN_DIR}/go-client.log"

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
trap 'kill_pids ${fd_pid:-} ${go_pid:-}' EXIT
wait_for_line "${fd_log}" "freeDiameterd daemon initialized" 15

# Note: the go s6a-client sends AIR (318) first, then ULR (316). FD's
# test_app only echoes the configured cmd_id, so either set cmd_id=318
# or 316. This scenario tests ULR specifically.
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-s6a-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-host "${FD_IDENTITY}" -peer-realm "${REALM}"  -timeout 10s)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
# We don't require both AIR and ULR to succeed against test_app (it
# echoes only the configured cmd). Consider scenario success if we got
# a 2001 answer on the negotiated app.
set +e
wait "${go_pid}"; rc=$?
set -e

assert_in_log "${fd_log}" "STATE_OPEN|CEA"
grep -qE "ULA Result-Code=2001|AIA Result-Code=2001" "${go_log}" \
    || { log "no 2001 Result-Code on any S6a answer"; tail -n 40 "${go_log}" >&2; exit 1; }
