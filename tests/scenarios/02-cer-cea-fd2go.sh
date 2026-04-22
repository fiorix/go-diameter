#!/bin/bash
# 02-cer-cea-fd2go: go-diameter server listens, freeDiameterd initiates
# CER, completes handshake.

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

go_log="${SCEN_DIR}/go-server.log"
fd_log="${SCEN_DIR}/fd.log"

# Start go server first so FD has something to connect to.
wrapper=$(ct_wrapper_path "${go_log}" "${CT_BIN}/diam-smoke-server"  -addr ":${GO_PORT}"  -origin-host "${GO_IDENTITY}" -origin-realm "${REALM}")
sdme exec "${GO_CT}" -- /bin/sh "${wrapper}" &
go_pid=$!
trap 'kill_pids ${go_pid:-} ${fd_pid:-}; stop_processes_in "${GO_CT}" diam-smoke-server' EXIT

wait_for_line "${go_log}" "listening on" 10 || { log "go server didn't start"; exit 1; }

wrapper=$(ct_wrapper_path "${fd_log}" /usr/local/bin/freeDiameterd -c "${ct_scen}/freeDiameter.conf" -d)
sdme exec "${FD_CT}" -- /bin/sh "${wrapper}" &
fd_pid=$!

# FD logs "STATE_OPEN" once CER/CEA completes.
wait_for_line "${fd_log}" "STATE_OPEN|CEA received" 15 \
    || { log "CER/CEA did not complete"; exit 1; }

assert_in_log "${fd_log}" "STATE_OPEN|CEA received"
