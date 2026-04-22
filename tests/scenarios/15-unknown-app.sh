#!/bin/bash
# 15-unknown-app: go CC client advertises an App-Id that neither side
# knows; FD's CEA should return Result-Code=5010 (NO_COMMON_APPLICATION)
# or drop the connection.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# FD base config advertises only base (App-Id 0). go will request
# App-Id 4 which FD doesn't know since we don't load dict_dcca_3gpp.
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

# diam-cc-client advertises Auth-Application-Id=4 in CER. FD hasn't
# loaded Credit-Control, so no common application.
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-cc-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-realm "${REALM}"  -app-id 4 -timeout 3s)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
wait "${go_pid}" 2>/dev/null || true

grep -qE "5010|NO_COMMON_APPLICATION|3002|UNABLE_TO_DELIVER" \
    "${fd_log}" "${go_log}" \
    || { log "expected 5010/3002 for unknown app"; \
         tail -n 40 "${fd_log}" >&2; tail -n 40 "${go_log}" >&2; exit 1; }
