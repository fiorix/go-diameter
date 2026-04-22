# Run #13 results (2026-04-22, arm64)

Baseline from the first working end-to-end run of the interop harness.

## Environment

- Host: Linux arm64 (aarch64), systemd-nspawn via sdme
- Base rootfs: `docker.io/ubuntu` → Ubuntu 24.04.4 LTS
- Go: 1.25.0 (downloaded during rootfs build, not the apt `golang-go` package)
- freeDiameter: open5gs/freeDiameter `main` branch built from source
- Disabled FD extensions (awkward deps): dict_json, app_acct, app_diameap,
  app_radgw, app_sip, app_redirect, dbg_interactive

## Summary

7 / 17 pass, 10 fail.

| # | Scenario | Status | Why |
|---|---|---|---|
| 01 | cer-cea-go2fd               | PASS | go client → FD, full CER/CEA/DPR |
| 02 | cer-cea-fd2go               | PASS | FD → go server, CER/CEA |
| 03 | dwr-dwa-go2fd               | PASS | go client watchdog visible in FD log |
| 04 | dwr-dwa-fd2go               | PASS | FD TwTimer triggers DWRs to go |
| 05 | cer-cea-reject              | PASS | FD returns 3010 DIAMETER_UNKNOWN_PEER |
| 06 | echo-go2fd                  | FAIL | routing: FD 3002 UNABLE_TO_DELIVER |
| 07 | echo-fd2go                  | FAIL | SIGUSR1 client-mode trigger not firing echoes |
| 08 | s6a-ulr-go2fd               | FAIL | 3001 COMMAND_UNSUPPORTED; test_app handles one cmd |
| 09 | s6a-ulr-fd2go               | FAIL | same SIGUSR1 flow as 07 |
| 10 | s6a-air-both-directions     | FAIL | combines 07 and 08 issues |
| 11 | ccr-ro-go2fd                | FAIL | 5001 AVP_UNSUPPORTED (FD lacks dict_dcca here) |
| 12 | ccr-ro-fd2go                | FAIL | SIGUSR1 flow (as 07) |
| 13 | ccr-gx-go2fd                | FAIL | go dict load conflict: CC(272) already registered |
| 14 | ccr-gx-fd2go                | FAIL | same go dict conflict as 13 |
| 15 | unknown-app                 | PASS | FD returns 3002 when app not negotiated |
| 16 | missing-mandatory-avp       | FAIL | 5001 AVP_UNSUPPORTED (same dict gap as 11) |
| 17 | invalid-avp-length          | PASS | FD rejects malformed AVP |

## What works end-to-end

- TCP transport over loopback in a single sdme pod.
- Self-signed cert + No_TLS peers (FD init quirk).
- CER/CEA handshake in both directions, with application-id negotiation.
- DWR/DWA watchdog in both directions.
- DPR/DPA clean disconnect (go side; FD sends DPA 2001 on receipt).
- Unknown peer and malformed AVP rejection paths.

## Root causes of the remaining failures

1. **Routing policy (06, 07, 09, 10, 12):** FD defaults to forwarding
   requests rather than dispatching locally when rt_default is not
   configured with a working conf file. Loading `rt_default.fdx` with
   an empty arg fails (`No such file or directory`); a minimal
   conf file that teaches it "local realm = us" is needed.

2. **test_app cmd_id is singular (08, 10):** The extension echoes only
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

5. **test_app client mode SIGUSR1 (07, 09, 10, 12):** Handshake succeeds
   but the subsequent signals don't produce answers in the FD log.
   Either the signal isn't reaching test_app, or test_app needs an
   additional routing piece. Verifying `pkill -USR1 -x freeDiameterd`
   inside the container against the test_app flow is the next step.

## What's next

- Ship a minimal `rt_default.conf` that routes local-realm requests to
  self (likely unblocks 06 and parts of 07, 09, 10).
- Try loading `dict_nasreq` then `dict_dcca` *before* `test_app` in
  scenarios 11-12, 16, to see if the "conflicting entry" resolves.
- For Gx (13, 14), register the Gx app programmatically in
  `tests/go/diam-cc-client` rather than loading the whole XML.
- Split S6a scenarios so each exercises a single command.
- Investigate the SIGUSR1 path (stdout buffering, test_app log level).

Each of these is a concrete, tractable follow-up — none requires
redesigning the harness.
