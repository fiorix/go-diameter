#!/bin/bash
# Top-level runner for go-diameter <-> FreeDiameter interop scenarios.
#
# Usage:
#   sudo ./tests/run.sh                # run all scenarios
#   sudo ./tests/run.sh 01-cer-cea-*   # run a subset (glob ok)
#   sudo ./tests/run.sh --keep-pod     # leave pod+containers up after
#
# Prereqs (one-time):
#   sudo sdme fs import ubuntu docker.io/ubuntu:24.04 -v --install-packages=yes
#   sudo sdme fs build godiam-fd tests/fs/build.conf -v

set -euo pipefail

TESTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tests/lib.sh
source "${TESTS_DIR}/lib.sh"

KEEP_POD=0
SCENARIOS=()

for arg in "$@"; do
    case "${arg}" in
        --keep-pod) KEEP_POD=1 ;;
        -h|--help)
            sed -n '2,12p' "$0" | sed 's/^# \{0,1\}//'
            exit 0 ;;
        *) SCENARIOS+=("${arg}") ;;
    esac
done

need_root
have_rootfs || die "rootfs '${ROOTFS}' not found; see prereqs in this script"

ensure_pod
ensure_container "${GO_CT}" "${GO_IDENTITY}"
ensure_container "${FD_CT}" "${FD_IDENTITY}"
ensure_certs
build_go_binaries

# Discover scenarios.
shopt -s nullglob
if (( ${#SCENARIOS[@]} == 0 )); then
    SCENARIO_FILES=("${TESTS_DIR}"/scenarios/*.sh)
else
    SCENARIO_FILES=()
    for pat in "${SCENARIOS[@]}"; do
        matches=("${TESTS_DIR}"/scenarios/${pat}*.sh)
        (( ${#matches[@]} > 0 )) || { log "no match for '${pat}'"; continue; }
        SCENARIO_FILES+=("${matches[@]}")
    done
fi

mkdir -p "${LOG_ROOT}"

pass=0; fail=0
declare -a results
for scen in "${SCENARIO_FILES[@]}"; do
    name="$(basename "${scen}" .sh)"
    dir="${REPO_DIR}/tests/logs/${name}"
    mkdir -p "${dir}"
    log "== scenario: ${name} =="
    start_ts=$(date +%s)
    if SCEN_NAME="${name}" SCEN_DIR="${dir}" bash "${scen}"; then
        status="PASS"; pass=$((pass+1))
    else
        status="FAIL"; fail=$((fail+1))
    fi
    dur=$(( $(date +%s) - start_ts ))
    results+=("${status}  ${name}  ${dur}s")
    cleanup_scenario || true
done

echo
echo "==== RESULTS ===="
for r in "${results[@]}"; do echo "${r}"; done
echo "-----------------"
echo "pass=${pass} fail=${fail}"

if (( KEEP_POD == 0 )); then
    log "stopping containers (use --keep-pod to leave them up)"
    sdme stop "${GO_CT}" "${FD_CT}" 2>/dev/null || true
fi

(( fail == 0 ))
