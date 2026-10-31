#!/bin/bash
# 04-dwr-dwa-fd2go: FD connects to go server; FD's TwTimer fires DWR
# after the connection is idle. Test uses a compressed timer via FD conf
# override.
#
# The go server runs with a fixed Origin-State-Id of 1, which never
# matches FD's (FD uses its start time). A CEA must carry the sender's
# own Origin-State-Id, the same one its DWA carries (RFC 6733 section
# 8.16). A CEA that echoes the peer's value disagrees with the DWA, and
# FD resets the connection on the first watchdog exchange. go-diameter
# did that from v4.5.0 until v4.6.0.

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
go_osid=1

wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-server"  -addr ":${GO_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}"  -origin-state-id "${go_osid}")
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
trap 'kill_pids ${go_pid:-} ${fd_pid:-}; stop_processes_in "${GO_CT}" diam-smoke-server' EXIT
wait_for_line "${go_log}" "listening on" 10
wait_for_port "${GO_PORT}" 10

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!
wait_for_line "${fd_log}" "STATE_OPEN|CEA received" 15

# Two DWAs take two TwTimer periods (6s each). By the second one, FD has
# finished processing the first, including any Origin-State-Id reset.
wait_for_count "${fd_log}" "'Device-Watchdog-Answer'" 2 25 || true

assert_in_log "${fd_log}" \
    "Capabilities-Exchange-Answer.*Origin-State-Id\(278\)\[-M\]=${go_osid} \(0x${go_osid}\)"
if grep -qE "Received a new Origin-State-Id|'STATE_CLOSED'" "${fd_log}"; then
    log "FD reset the connection over a changed Origin-State-Id"
    grep -E "Origin-State-Id\(278\)|Received a new Origin-State-Id|STATE_" "${fd_log}" >&2 || true
    exit 1
fi
dwa_count=$(grep -cE "'Device-Watchdog-Answer'" "${fd_log}" || true)
[[ "${dwa_count}" -ge 2 ]] || { log "expected >=2 DWAs, got ${dwa_count}"; exit 1; }
