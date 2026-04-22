#!/bin/bash
# 04-dwr-dwa-fd2go: FD connects to go server; FD's TwTimer fires DWR
# after the connection is idle for 30s. Test uses a compressed timer
# via FD conf override.

set -euo pipefail
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source "${TESTS_DIR}/lib.sh"
: "${SCEN_DIR:?}"
ct_scen="${CT_REPO}/tests/logs/$(basename "${SCEN_DIR}")"

# Render base conf and append TcTimer / TwTimer overrides.
FD_IDENTITY="${FD_IDENTITY}" FD_REALM="${REALM}" FD_PORT="${FD_PORT}" \
FD_PEER_IDENTITY="${GO_IDENTITY}" FD_PEER_PORT="${GO_PORT}" \
FD_CERT_DIR="${CT_CERT}" FD_EXT_DIR="${FD_EXT_DIR}" \
    envsubst < "${TESTS_DIR}/configs/fd-base.conf.tmpl" \
    > "${SCEN_DIR}/freeDiameter.conf"
cat >> "${SCEN_DIR}/freeDiameter.conf" <<EOF
TcTimer = 3;
TwTimer = 6;
EOF

go_log="${SCEN_DIR}/go-server.log"
fd_log="${SCEN_DIR}/fd.log"

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-server"  -addr ":${GO_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}")
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
trap 'kill_pids ${go_pid:-} ${fd_pid:-}; stop_processes_in "${GO_CT}" diam-smoke-server' EXIT
wait_for_line "${go_log}" "listening on" 10

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
wait_for_line "${fd_log}" "STATE_OPEN|CEA received" 15
sleep 15  # let TwTimer (6s) fire at least twice

dwr_count=$(grep -cE "Device-Watchdog-Request|DWR" "${fd_log}" || true)
[[ "${dwr_count}" -ge 2 ]] || { log "expected >=2 DWRs, got ${dwr_count}"; exit 1; }
