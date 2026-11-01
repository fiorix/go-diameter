# go-diameter interop tests

End-to-end scenarios that exercise go-diameter against the
open5gs/freeDiameter fork. Both peers run in sdme (systemd-nspawn)
containers sharing a pod network namespace; communication is over
TCP/loopback.

## Prerequisites

Host:

- Linux with systemd >= 255
- root (sdme needs it)
- `sdme` in PATH (https://github.com/fiorix/sdme)
- `gettext-base` for `envsubst` (on Debian/Ubuntu: `apt install gettext-base`)
- `openssl` for the one-time self-signed cert generation
- `nsenter` (util-linux) and `ss` (iproute2), used from the host to add
  an address to the pod netns and to wait for listening ports

One-time rootfs build:

```
sudo sdme fs import ubuntu docker.io/ubuntu -v --install-packages=yes
sudo sdme fs build godiam-fd tests/fs/build.conf -v
```

This builds the freeDiameter daemon + extensions into the rootfs,
installs the Go toolchain (1.27.3 by default, `GO_VERSION` in
`tests/fs/scripts/install-fd.sh`) and fills the Go module cache from
`go.mod`, since the test pod has no network. Subsequent scenario runs
reuse the rootfs. Rebuild it with `-f` when `go.mod` or `go.sum`
change. Last built on Ubuntu 26.04.

## Running

```
sudo ./tests/run.sh                 # all scenarios
sudo ./tests/run.sh 01-cer-cea-go2fd 03-dwr-dwa-go2fd
sudo ./tests/run.sh --keep-pod      # leave containers up for debugging
```

Per-scenario logs land in `tests/logs/<scenario>/`:

- `freeDiameter.conf` — rendered FD config used for that run
- `fd.log` — freeDiameterd stderr (includes `dbg_msg_dumps` output)
- `go-*.log` — go test binary output

## Topology

```
+-------------- pod: gdiam-interop -----------------+
|                                                   |
|  container "godiam"        container "fdiam"      |
|  hostname godiam.test.local   fdiam.test.local    |
|  port 3868 (go server)     port 3869 (FD)         |
|  /repo bind-mounted        /repo bind-mounted     |
|                                                   |
+--- loopback 127.0.0.1 (shared netns) -------------+
```

Both containers run from the single `godiam-fd` rootfs; they differ
only in `/etc/hostname` (bind-mounted per container).

## Scenarios

Numbered files under `tests/scenarios/`. Each script renders its own
freeDiameter config from `tests/configs/*.tmpl`, starts the peers,
asserts against logs, and returns 0/non-zero.

Status on the latest run; see `RESULTS.md` for per-failure diagnosis.
Three consecutive runs give the same result.

| # | Name | Status | Coverage |
|---|---|---|---|
| 01 | cer-cea-go2fd            | PASS | CER/CEA + DPR |
| 02 | cer-cea-fd2go            | PASS | CER/CEA (FD initiates) |
| 03 | dwr-dwa-go2fd            | PASS | Watchdog, go-side timer |
| 04 | dwr-dwa-fd2go            | PASS | Watchdog, FD TwTimer, Origin-State-Id in CEA and DWA |
| 05 | cer-cea-reject           | PASS | 3010 DIAMETER_UNKNOWN_PEER |
| 06 | echo-go2fd               | FAIL | rt_default missing -> 3002 |
| 07 | echo-fd2go               | FAIL | SIGUSR1 path for test_app client |
| 08 | s6a-ulr-go2fd            | FAIL | test_app handles single cmd |
| 09 | s6a-ulr-fd2go            | FAIL | SIGUSR1 path |
| 10 | s6a-air-both-directions  | FAIL | go2.test.local rejected with 3010 |
| 11 | ccr-ro-go2fd             | FAIL | 5001 without dict_dcca |
| 12 | ccr-ro-fd2go             | FAIL | SIGUSR1 path |
| 13 | ccr-gx-go2fd             | FAIL | Gx XML re-registers CC(272) |
| 14 | ccr-gx-fd2go             | FAIL | same go dict conflict |
| 15 | unknown-app              | PASS | 3002 on CCR for un-negotiated app |
| 16 | missing-mandatory-avp    | FAIL | 5001 without dict_dcca |
| 17 | invalid-avp-length       | PASS | FD rejects malformed AVP |

## Scope & limits

- **TCP only**: SCTP in nspawn requires extra kernel/privilege setup.
- **No TLS/IPsec**: FD still needs a `TLS_Cred` directive to init,
  satisfied by a throwaway self-signed cert under `configs/certs/`.
- **S6a / Ro / Gx semantic depth is limited** because pure FD has no
  HSS/PCRF logic; the `test_app.fdx` extension is used as a
  parameterizable echo for the configured App-Id/Cmd-Code. This
  validates CER application negotiation, command framing, and wire
  format compatibility. Semantic correctness of 3GPP flows is covered
  by go-diameter's own unit tests (`diam/sm/s6a_client_server_test.go`).
- **Error scenarios 16-17** rely on a small custom binary,
  `tests/go/diam-bad-client`, that bypasses the normal marshaller.

## Adding a new scenario

1. Add a new numbered script under `tests/scenarios/`.
2. Render one of the `configs/*.tmpl` files with `envsubst`.
3. Start FD and Go binaries via `sdme exec`. Before dialing a peer,
   `wait_for_port` on its port: both peers log that they are up before
   their socket is bound.
4. Assert with `grep`/`assert_in_log`.
5. Make it executable; `run.sh` auto-discovers it.

## Debugging a failing scenario

```
sudo ./tests/run.sh 07-echo-fd2go --keep-pod
sudo sdme join godiam   # shell inside the go-diameter container
sudo sdme join fdiam    # shell inside the FD container
cat tests/logs/07-echo-fd2go/fd.log
```

Enable `dbg_msg_dumps` at a higher level (e.g. `0xFFFF`) in
`tests/configs/fd-base.conf.tmpl` for more verbose FD decode output.
