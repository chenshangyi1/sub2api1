## 2026-10-07 - Task: SiliconFlow DeepSeek V3.2 base pricing and source synchronization

### What was done

Corrected V3.2 site base prices to 4/6/0.40 per million tokens, protected them from dynamic USD catalog and generic Flash rewrites, preserved explicit pricing overrides and effective multipliers. Audited other configured DeepSeek models without changing their prices. User selected sub2api1/main for GitHub synchronization. Server synchronization is source-only to its existing matching staging tree; no runtime deployment is authorized.

### Testing

New regression tests reproduced the wrong Flash rate before the fix. After the fix all four V3.2 regression tests passed using the package's production GoFiles and focused test file, with GOEXPERIMENT=jsonv2. The server command built successfully. Full unit package compilation is blocked by preexisting duplicate definitions and stale account-stats test APIs. Expanded tests exposed underscore-prefix normalization and V4 date-suffix catalog matching failures; the same failures were reproduced with the unchanged 847ce5afa billing source through a Go overlay. See the audit note. Runtime preservation is checked during source synchronization.

### Notes

Changed billing_service.go, moved/expanded the V3.2 tests out of deepseek_pricing_test.go into deepseek_v32_pricing_test.go, and added the audit/handoff/progress notes. Local rollback is a revert of this operation's commit; server source rollback uses its operation backup. No historical balances/bills, group settings, compose files or running images were changed.

## 2026-10-07 - Task: Local gateway verification and candidate-image preparation

### What was done

The user authorized building a new image after local tests, while retaining the running production container. Added three V3.2 regressions for the actual gateway token-cost path, zero/cache-only requests and concurrent partial/free group overrides. Added a minimal candidate-image recipe and a Dockerfile-specific context filter. Compared all backend Go/module source files and all 221 frontend files with the running image; runtime-source differences are confined to this pricing change.

### Testing

All seven focused V3.2 tests passed, including three repetitions with the Go race detector. The real-bill token tuple 23/137/4352 produces base cost 0.0026548 and actual cost 0.00238932 at multiplier 0.9. Expanded billing/tier/gateway regressions reported 210 passing test/subtest events and two Grok 4.6 failures. Both Grok failures were reproduced against unchanged 847ce5afa source with an overlay and identical mismatches, confirming they predate this change. Full-package compilation limitations and previously verified V4 baseline failures remain unresolved; no claim of a globally green suite.

### Notes

Changed the focused test, DEPLOYMENT_HANDOFF.md and added deploy/Dockerfile.deepseek-v32-candidate plus its .dockerignore. A candidate will replace only /app/sub2api over the verified current runtime image; construction and isolated version checks do not deploy it. Test evidence is in ignored work files. Source rollback is a revert of this task’s commits; existing server source backup remains /opt/sub2api/backups/v32-source-20261007T081207Z. The unrelated api_key file remains untouched.

## 2026-10-07 - Task: Candidate image built and isolated Linux checks passed

### What was done

Built sub2api1:deepseek-v32-cny-9c3c6fbd8-20261007 from the canonical server build root, over the exact current runtime image. Replaced only the executable with the locally validated linux/amd64 artifact from source 9c3c6fbd809441d657da26e9039f29ca8c3ab009. Recorded image/executable hashes, version and backups in DEPLOYMENT_HANDOFF.md. No rollout occurred.

### Testing

Candidate image ID e8f9acf6dfcab55ace34151be7526456acdb0f2706592a3c9c4e4a44e9b18f9a; executable SHA256 b8cf200a71a28313d7978d8da0b3953756d5e14cbb9e925ce1bffa67e64d9f9b. Isolated version and checksum checks passed. All seven pricing regressions also passed as a cross-compiled Linux test executable inside the read-only, network-disabled candidate, without production volumes or credentials. Runtime image configuration matched the base, and the existing production container's identity/start time/image remained unchanged with readiness HTTP 200. This supplements the local/race tests; unrelated baseline failures remain as previously documented.

### Notes

This step changes only handoff/progress documentation locally; uploaded binaries and test evidence reside under the existing server build root's outputs/deepseek-v32-candidate. The unused image can be discarded without runtime rollback. Source sync backup: /opt/sub2api/backups/v32-candidate-20261007T082356Z. Old /app Go source is inherited in the binary-overlay image, so future compilation must use the recorded repository revision. No service restart or production data write was performed.

## 2026-10-07 - Task: Passive CF/request latency evidence collection

### What was done

Collected bounded live internal HTTP timing metadata in memory, without retaining raw packets, credentials, cookies, conversation contents or user IPs. Two natural SiliconFlow DeepSeek V3.2 successes were joined by response client/request IDs to CF-Ray, application logs and read-only usage records. Both entered CF LHR; origin first SSE events were 3.725/4.825 seconds, application totals 35.771/53.383 seconds, maximum internal SSE gaps 0.728/2.069 seconds. Historical application/Caddy logs and database request metadata do not retain CF-Ray, so earlier slow requests' entry nodes cannot be reconstructed from these records. The main historical delays occurred within upstream forwarding; provider queuing, transport, historical stream stalls and downstream cancellation remain distinguishable hypotheses, not proven causes.

### Verification and limitations

Capture windows were 19:13:04–19:17:16 and 19:17:28–19:24:28 UTC+8, 4009 parsed packets total, zero kernel drops. A keepalive idle-time accounting issue in the initial temporary parser was detected using unauthenticated 401 probes, fixed and checked with a synthetic reuse case. Initial sample timing was recomputed from its matched application access log; raw initial start/duration fields are not relied upon. Whitelisted synthetic framing/SSE/DONE timing and sensitive-field exclusion passed. No model calls were made by the agent. Two normal live samples do not exclude intermittent CF/downstream or server-to-provider issues, and origin-to-Caddy capture does not prove end-user delivery.

### Notes

Scripts, redacted metadata and report are in ignored work/; see work/latency-investigation-report.md and work/live-latency-joined.json. No Caddy/DNS/database/Compose/source/image change or service restart occurred. All temporary capture processes exited. Final container ID/image/StartedAt remain unchanged, healthy, restart_count=0, readyz HTTP 200. Candidate pricing image remains unused. No server rollback is necessary; only temporary local diagnostic artifacts and this progress entry were added. No GitHub push was performed for diagnostic artifacts.

## 2026-10-07 - Task: Deploy verified V3.2 pricing candidate

### What was done

After the user explicitly authorized deployment, created and validated a root-only PostgreSQL backup and runtime snapshots at `/opt/sub2api/backups/v32-rollout-20261007T161509Z`. The candidate's read-only migration check passed using the current config and environment. Changed only the `sub2api` image field in the live Compose file and recreated only that app service; PostgreSQL, Redis, persistent volumes, frontend, Caddy and existing group/user multipliers were preserved.

### Verification

The candidate `sub2api1:deepseek-v32-cny-9c3c6fbd8-20261007` is healthy with restart count 0 and the expected image/executable hashes. Origin and public `/readyz` returned HTTP 200. Public model-plaza shows DeepSeek V3.2 base prices of 4/6/0.40 per million tokens and group 16 rate 0.9; unrelated model prices and group rates are unchanged. PostgreSQL and Redis container IDs/start times, runtime mounts/settings, frontend symlink and Caddy/config checksums are unchanged. Startup logs had no ERROR/FATAL/PANIC entries. The readiness timeline measured one unavailable sample during the single-container switch; healthy state was reached about seven seconds after switching.

### Rollback

Restore `/opt/sub2api/backups/v32-rollout-20261007T161509Z/docker-compose.local.yml` and run the app-only Compose recreate. Keep the database backup for data recovery only; do not restore it for an image rollback.

## 2026-10-11 - Task: Deploy DeepSeek V4 Pro pricing candidate

The user authorized replacing the verified candidate image. The candidate's only pending migration, `256_account_stats_rules_user_billing.sql`, was applied successfully as an additive `BOOLEAN NOT NULL DEFAULT FALSE` column before the app switch. The PostgreSQL dump and runtime backup are retained at `/opt/sub2api/backups/deepseek-v4-pro-rollout-20261010T175122Z`.

Production now runs `sub2api1:deepseek-v4-pro-11f9760e-20261011` with image ID `sha256:0d136bc6b2216825cf4431bcb3914b49e9ca21307697af6126626c647a770b73`, container `ad9ff96cc91660356d62007bf7c2ad9acaa81fbc4a4951b14f2019151582ebb4`, and executable SHA256 `47bd98a37359060c96a7c66a41d5cf8044bb9a0ce395215398981e147b58e461`. The app is healthy, restart count 0, origin/public `/readyz` are HTTP 200, PostgreSQL and Redis were not restarted, and startup logs contain no ERROR/FATAL/PANIC entries. Roll back only the app Compose image from the recorded backup; do not restore the database for image rollback.

## 2026-10-11 - Task: Retain current/previous images and remove old tags

Kept the running main-app image plus the previous deployed gateway-cors image and the same-day account-rule candidate. Kept upstream-ops current `jws-browser-20261009` plus previous `no-profit-20261009`, along with the active PostgreSQL and Redis images. Removed four unreferenced historical tags: DeepSeek V3.2, Grok build-fast, upstream-ops profit-reconcile, and upstream-ops sub2api-sync-home. No containers, volumes or build cache were pruned. Before/after inventories and deletion checks are recorded in `/opt/sub2api/backups/image-cleanup-20261010T180533Z`; all services remained healthy with restart count zero.


## 2026-10-11 - Task: Verify Claude JWS source and build unused candidate

Synchronized compiled billing source 2b4f5d51a926256f2c7727b42f98b2e301f0e273 to sub2api1/main and the existing server staging tree, after backups and source hash verification (3,427 files match). Built sub2api1:claude-jws-2b4f5d51a-20261011 with image ID sha256:017b5b91dd650e44e9cc395cd5edfeb91e4e9eab9d13696e651daec0edc519c1 and executable SHA256 f809796957a4fb1b089efb56ec3fae949656983b69a57c0e06d6213ebba80717.

Financial regression 496, race 64, isolated candidate 64 test/subtest events, focused handler checks, frontend plaza 25 and main Vue typecheck pass. Full backend and frontend failures plus lint/dev typecheck blockers are documented; no full-suite success claim. Final documentation and legacy-compatible COPY recipe corrections do not alter compiled billing source.

Evidence: /opt/sub2api/build_context/outputs/claude-jws-candidate; original source backup: /opt/sub2api/backups/claude-jws-source-20261011T135045Z. All four running containers and Compose/Caddy/frontend pointers are unchanged. Candidate is unused. Nine successful-build readiness samples and final origin/public checks returned 200. No restart, live image switch, migration, DB write or paid model request.
