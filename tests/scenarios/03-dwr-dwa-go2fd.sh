#!/bin/bash
# 03-dwr-dwa-go2fd: go client connects, enables its watchdog, holds
# the connection open long enough for several DWR/DWA round-trips.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

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

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-client"  -addr "127.0.0.1:${FD_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -peer-host "${FD_IDENTITY}" -peer-realm "${REALM}"  -watchdog 2s -timeout 7s -dpr)
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
wait "${go_pid}"

assert_in_log "${go_log}" "CER/CEA handshake OK"
# FD's dbg_msg_dumps prints "Device-Watchdog-Request" when it arrives.
# At least 2 DWR round-trips expected in 7s with 2s interval.
dwr_count=$(grep -cE "Device-Watchdog-Request|DWR" "${fd_log}" || true)
[[ "${dwr_count}" -ge 2 ]] || { log "expected >=2 DWRs, got ${dwr_count}"; exit 1; }
