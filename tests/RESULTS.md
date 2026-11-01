# Results (2026-11-01, x86_64)

go-diameter main after v4.6.0 (df25c21). Three consecutive full runs
gave the same status for every scenario, about 105 s per run.

## Environment

- Host: Linux x86_64, systemd-nspawn via sdme 0.21.1
- Base rootfs: `docker.io/ubuntu` -> Ubuntu 26.04.1 LTS
- Go: 1.27.3 (downloaded during rootfs build, not the apt `golang-go`
  package); the module cache is filled during the build too, since the
  pod has no network
- freeDiameter: open5gs/freeDiameter `main` branch built from source
  with CMake 4 (`CMAKE_POLICY_VERSION_MINIMUM=3.5`, CMP0022 set to NEW)
- Disabled FD extensions (awkward deps): dict_json, app_acct, app_diameap,
  app_radgw, app_sip, app_redirect, dbg_interactive
- Pod netns: `lo` plus `dummy0` with 192.0.2.1/32, because freeDiameterd
  refuses to start without a non-loopback address

## Summary

7 / 17 pass, 10 fail, on each of three runs.

| # | Scenario | Status | Why |
|---|---|---|---|
| 01 | cer-cea-go2fd               | PASS | go client -> FD, full CER/CEA/DPR |
| 02 | cer-cea-fd2go               | PASS | FD -> go server, CER/CEA |
| 03 | dwr-dwa-go2fd               | PASS | go client watchdog visible in FD log |
| 04 | dwr-dwa-fd2go               | PASS | FD TwTimer DWRs answered; CEA and DWA carry the same Origin-State-Id |
| 05 | cer-cea-reject              | PASS | FD returns 3010 DIAMETER_UNKNOWN_PEER |
| 06 | echo-go2fd                  | FAIL | routing: FD 3002 UNABLE_TO_DELIVER |
| 07 | echo-fd2go                  | FAIL | SIGUSR1 client-mode trigger produces no requests |
| 08 | s6a-ulr-go2fd               | FAIL | 3001 COMMAND_UNSUPPORTED; test_app handles one cmd |
| 09 | s6a-ulr-fd2go               | FAIL | same SIGUSR1 flow as 07 |
| 10 | s6a-air-both-directions     | FAIL | FD -> go answered; go2.test.local rejected with 3010 |
| 11 | ccr-ro-go2fd                | FAIL | 5001 AVP_UNSUPPORTED (FD lacks dict_dcca here) |
| 12 | ccr-ro-fd2go                | FAIL | SIGUSR1 flow (as 07) |
| 13 | ccr-gx-go2fd                | FAIL | go dict load conflict: CC(272) already registered |
| 14 | ccr-gx-fd2go                | FAIL | same go dict conflict as 13; the cc server never starts |
| 15 | unknown-app                 | PASS | FD returns 3002 when app not negotiated |
| 16 | missing-mandatory-avp       | FAIL | 5001 AVP_UNSUPPORTED (same dict gap as 11) |
| 17 | invalid-avp-length          | PASS | FD rejects malformed AVP |

Scenario 04 is the regression test for the Origin-State-Id fix in
v4.6.0. The go server runs with Origin-State-Id 1, which never matches
freeDiameter's. Against v4.5.1 library code the CEA carries
freeDiameter's own value back, freeDiameter logs "Received a new
Origin-State-Id ... resetting the connection" on the first DWA, and the
scenario fails.

## What works end-to-end

- TCP transport over loopback in a single sdme pod.
- Self-signed cert + No_TLS peers (FD init quirk).
- CER/CEA handshake in both directions, with application-id negotiation.
- DWR/DWA watchdog in both directions, with a consistent Origin-State-Id.
- DPR/DPA clean disconnect (go side; FD sends DPA 2001 on receipt).
- Unknown peer and malformed AVP rejection paths.

## Harness reliability

Earlier runs varied from run to run. Three causes were fixed:

- Peers log "listening on" or "daemon initialized" before their socket
  is bound, so the dialing side sometimes got "connection refused"
  (05, 08, 09, 10 and 12). Scenarios now wait for the port itself
  (`wait_for_port`), and run.sh waits for both ports to be free before
  each scenario.
- run.sh reused `tests/logs/<scenario>/`, so a scenario could match the
  previous run's log before the new peer truncated it. It now starts
  each scenario with an empty directory.
- 04 depended on whether both peers started in the same second, since
  both seeded Origin-State-Id from their start time.

## Root causes of the remaining failures

1. **Routing policy (06, 07, 09, 12):** FD defaults to forwarding
   requests rather than dispatching locally when rt_default is not
   configured with a working conf file. Loading `rt_default.fdx` with
   an empty arg fails (`No such file or directory`); a minimal
   conf file that teaches it "local realm = us" is needed.

2. **test_app cmd_id is singular (08):** The extension echoes only
   one (vendor, app, cmd) tuple per instance. S6a scenarios exercise
   both AIR (318) and ULR (316); at best we can validate one per
   scenario. Worth either splitting the S6a scenarios or sending a
   single command per run.

3. **dict_dcca vs test_app conflict at app-id 4 (11, 16):** Loading
   both causes "Conflicting entry in the dictionary" for application
   id 4. When dict_dcca is skipped (to avoid the conflict), FD doesn't
   know CC-Request-Type (416) etc., and rejects the go client's CCR
   with 5001 AVP_UNSUPPORTED. Resolution likely involves loading
   dict_nasreq + dict_dcca *before* test_app so test_app uses the
   existing app object.

4. **go dict self-conflict for Gx (13, 14):** Loading
   `diam/dict/testdata/gx_credit_control.xml` on top of the default
   dict re-registers Credit-Control (cmd 272). Error: "cannot be added:
   index exists". The go-diameter default dict already has CC; loading
   Gx's XML wholesale duplicates it. A targeted Gx-only delta XML,
   or programmatic registration of just the Gx app mapping, is needed.
   In 14 the cc server exits on this error, so the scenario times out
   waiting for it to listen.

5. **test_app client mode SIGUSR1 (07, 09, 12):** Handshake succeeds
   but the subsequent signals don't produce answers in the FD log; the
   two 2001 lines the scenarios count are the CER/CEA exchange. Either
   the signal isn't reaching test_app, or test_app needs an additional
   routing piece.

6. **Second go identity not a known peer (10):** The FD -> go half of
   10 gets its 2001 answers. The go -> FD half dials with Origin-Host
   `go2.test.local`, which is not a configured peer, so FD answers the
   CER with 3010 DIAMETER_UNKNOWN_PEER. FD needs to accept that peer, or
   the go client needs to reuse the configured identity on a second
   connection.

## What's next

- Ship a minimal `rt_default.conf` that routes local-realm requests to
  self (likely unblocks 06 and parts of 07, 09, 12).
- Try loading `dict_nasreq` then `dict_dcca` *before* `test_app` in
  scenarios 11-12, 16, to see if the "conflicting entry" resolves.
- For Gx (13, 14), register the Gx app programmatically in
  `tests/go/diam-cc-client` and `tests/go/diam-cc-server` rather than
  loading the whole XML.
- Split S6a scenarios so each exercises a single command.
- Let FD accept `go2.test.local` in 10, for example with a second
  ConnectPeer entry or an acl_wl rule.
- Investigate the SIGUSR1 path (stdout buffering, test_app log level).
