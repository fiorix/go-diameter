#!/bin/bash
# Shared helpers for go-diameter <-> FreeDiameter interop scenarios.
# Sourced by tests/run.sh and tests/scenarios/*.sh.

set -euo pipefail

# Absolute path to the tests/ directory (the one that sourced us).
TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd -- "${TESTS_DIR}/.." && pwd)"

POD_NAME="${POD_NAME:-gdiam-interop}"
ROOTFS="${ROOTFS:-godiam-fd}"
GO_CT="${GO_CT:-godiam}"
FD_CT="${FD_CT:-fdiam}"

GO_IDENTITY="${GO_IDENTITY:-godiam.test.local}"
FD_IDENTITY="${FD_IDENTITY:-fdiam.test.local}"
REALM="${REALM:-test.local}"

GO_PORT="${GO_PORT:-3868}"
FD_PORT="${FD_PORT:-3869}"

# Directories inside the container (bind-mounted from repo).
CT_REPO=/repo
CT_BIN=/repo/tests/bin
CT_DICT=/repo/tests/dicts
CT_CERT=/repo/tests/configs/certs
FD_EXT_DIR="${FD_EXT_DIR:-/usr/local/lib/freeDiameter}"

LOG_ROOT="${LOG_ROOT:-/tmp/go-diameter-interop}"

log() { printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*" >&2; }
die() { log "FATAL: $*"; exit 1; }

need_root() {
    [[ ${EUID} -eq 0 ]] || die "run as root (sdme requires it)"
}

have_rootfs() {
    sdme fs ls 2>/dev/null | awk '{print $1}' | grep -Fxq "${ROOTFS}"
}

have_container() {
    sdme ps 2>/dev/null | awk 'NR>1 {print $1}' | grep -Fxq "$1"
}

ensure_pod() {
    if ! sdme pod ls 2>/dev/null | awk 'NR>1 {print $1}' | grep -Fxq "${POD_NAME}"; then
        log "creating pod ${POD_NAME}"
        sdme pod new "${POD_NAME}"
    fi
}

ensure_container() {
    local name="$1" hostname="$2"
    local hostfile="${LOG_ROOT}/hostname-${name}"
    local hostsfile="${LOG_ROOT}/hosts-${name}"
    mkdir -p "${LOG_ROOT}"
    printf '%s\n' "${hostname}" > "${hostfile}"
    cat > "${hostsfile}" <<EOF
127.0.0.1 localhost ${hostname}
::1       localhost ${hostname}
EOF
    if have_container "${name}"; then
        if ! sdme ps 2>/dev/null | awk -v n="${name}" '$1==n {print $2}' | grep -q running; then
            log "starting existing container ${name}"
            sdme start "${name}"
        fi
        return
    fi
    log "creating container ${name} (hostname=${hostname})"
    sdme create "${name}" \
        -r "${ROOTFS}" \
        --pod "${POD_NAME}" \
        -b "${REPO_DIR}:${CT_REPO}" \
        -b "${hostfile}:/etc/hostname:ro" \
        -b "${hostsfile}:/etc/hosts:ro"
    sdme start "${name}"
}

ct_exec() {
    local ct="$1"; shift
    sdme exec "${ct}" -- "$@"
}

ct_exec_bg() {
    local ct="$1" outfile="$2"; shift 2
    sdme exec "${ct}" -- "$@" >"${outfile}" 2>&1 &
    echo $!
}

# ct_wrapper_path <host_log> <cmd> [args...]
# Write a tiny wrapper script next to host_log that execs CMD with
# stdout+stderr redirected to the container-visible equivalent of
# host_log (via /repo bind mount). Prints the wrapper's container path
# on stdout; the caller invokes it with sdme exec and backgrounds
# directly so `wait` works on the resulting PID.
ct_wrapper_path() {
    local host_log="$1"; shift
    local ct_log="${CT_REPO}${host_log#${REPO_DIR}}"
    local wrapper="${host_log}.launch.sh"
    {
        printf '#!/bin/sh\nexec '
        for a in "$@"; do printf '%q ' "$a"; done
        printf "> %q 2>&1\n" "${ct_log}"
    } > "${wrapper}"
    chmod +x "${wrapper}"
    echo "${CT_REPO}${wrapper#${REPO_DIR}}"
}

build_go_binaries() {
    log "building Go test binaries in ${GO_CT}"
    ct_exec "${GO_CT}" /bin/sh -c "cd ${CT_REPO} && \
        mkdir -p ${CT_BIN} && \
        GOTOOLCHAIN=local \
        /usr/local/bin/go build -buildvcs=false -trimpath \
            -o ${CT_BIN}/ ./tests/go/..."
}

ensure_certs() {
    if [[ ! -s "${TESTS_DIR}/configs/certs/cert.pem" ]]; then
        log "generating self-signed FD certs"
        bash "${TESTS_DIR}/configs/certs/gen-certs.sh" "${TESTS_DIR}/configs/certs" fd-test
    fi
}

# render_fd_conf <template> <scenario_dir>
render_fd_conf() {
    local tmpl="$1" dest_dir="$2"
    mkdir -p "${dest_dir}"
    envsubst < "${tmpl}" > "${dest_dir}/freeDiameter.conf"
}

# render_testapp_conf <scenario_dir>
render_testapp_conf() {
    local dest_dir="$1"
    envsubst < "${TESTS_DIR}/configs/test_app.conf.tmpl" > "${dest_dir}/test_app.conf"
}

# start_fd <scenario_dir>
start_fd() {
    local dir="$1"
    local conf="${CT_REPO}/tests/logs/$(basename "${dir}")/freeDiameter.conf"
    local logf="${dir}/fd.log"
    log "starting freeDiameterd (conf=${conf})"
    sdme exec "${FD_CT}" -- sh -c \
        "exec freeDiameterd -c ${conf} -d 2>&1" \
        >"${logf}" 2>&1 &
    echo $!
}

# start_go <scenario_dir> <binary> [args...]
start_go() {
    local dir="$1" bin="$2"; shift 2
    local logf="${dir}/go-${bin}.log"
    log "starting ${bin} $*"
    sdme exec "${GO_CT}" -- "${CT_BIN}/${bin}" "$@" \
        >"${logf}" 2>&1 &
    echo $!
}

wait_for_line() {
    local file="$1" pattern="$2" timeout="${3:-10}"
    local elapsed=0
    while (( elapsed < timeout )); do
        if [[ -s "${file}" ]] && grep -qE -- "${pattern}" "${file}" 2>/dev/null; then
            return 0
        fi
        sleep 1
        elapsed=$(( elapsed + 1 ))
    done
    log "timeout waiting for: ${pattern} in ${file}"
    return 1
}

assert_in_log() {
    local file="$1" pattern="$2"
    if ! grep -qE -- "${pattern}" "${file}" 2>/dev/null; then
        log "FAIL: expected '${pattern}' in ${file}"
        log "---- last 40 lines of ${file} ----"
        tail -n 40 "${file}" >&2 || true
        return 1
    fi
}

kill_pids() {
    # First ensure the in-container long-running processes are dead,
    # otherwise sdme exec (machinectl shell) won't exit and the waits
    # below hang.
    sdme exec "${FD_CT:-fdiam}" -- /usr/bin/pkill -KILL -f freeDiameterd 2>/dev/null || true
    sdme exec "${GO_CT:-godiam}" -- /usr/bin/pkill -KILL -f "diam-" 2>/dev/null || true
    for pid in "$@"; do
        [[ -n "${pid:-}" ]] || continue
        kill -KILL "${pid}" 2>/dev/null || true
    done
    wait "$@" 2>/dev/null || true
}

stop_processes_in() {
    local ct="$1" pattern="$2"
    sdme exec "${ct}" -- /usr/bin/pkill -KILL -f "${pattern}" 2>/dev/null || true
}

cleanup_scenario() {
    stop_processes_in "${FD_CT}" freeDiameterd
    stop_processes_in "${GO_CT}" "diam-.*"
}
