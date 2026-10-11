# xinmc pricing candidate handoff — 2026-10-07

## Claude JWS candidate — 2026-10-11, built but NOT deployed

User authorized verification, sub2api1/main synchronization, server source staging synchronization and candidate construction, explicitly prohibiting a live image switch or service interruption. Candidate `sub2api1:claude-jws-2b4f5d51a-20261011`, image ID `sha256:017b5b91dd650e44e9cc395cd5edfeb91e4e9eab9d13696e651daec0edc519c1`, compiled source `2b4f5d51a926256f2c7727b42f98b2e301f0e273`, executable SHA256 `f809796957a4fb1b089efb56ec3fae949656983b69a57c0e06d6213ebba80717`. Final notes/recipe corrections are documentation/packaging only; the executable remains compiled from that source commit.

Claude uses JWS supplier cards in the candidate; Fable retail 1.2x and Max 3x surcharges are removed per the final instruction, while explicit cards and separately configured group/user multipliers remain authoritative. Laogou Grok prices are unchanged. Actual JWS evidence covers Opus 5.5, Sonnet 5 and Haiku 4.5; Fable write prices and some context/model/speed variants lack actual usage proof. Opus Fast retains its preexisting 2x rule; Standard is the verified JWS mode.

Local financial regression 496 and race 64 test/subtest events pass; handler-focused and frontend plaza 25 tests pass, main Vue typecheck passes. Candidate isolated Linux checks pass 64 selected test/subtest events, version/hash and runtime-config comparisons. Full backend/frontend/lint failures remain documented in `docs/claude-jws-candidate-20261011.md`; do not claim the full suite is green.

Artifacts/evidence: `/opt/sub2api/build_context/outputs/claude-jws-candidate`. Source staging: `/opt/sub2api/build_context/billing-847ce5afa-full`; original sync backup: `/opt/sub2api/backups/claude-jws-source-20261011T135045Z`. Canonical source local changes are preserved. Packaging uses a minimal context and command-local DOCKER_BUILDKIT=0 over the pinned current image; inherited image Go source is old, so future compilation must use the recorded Git revision.

All four production container IDs/images/start times/restart counts and health remain unchanged; the app is still deepseek-v4-pro-11f9760e-20261011, healthy, restart 0. Compose/Caddy hashes and frontend pointer are unchanged. Nine observed origin readiness samples and final origin/public checks returned 200. No deployment, migration, production data/card/multiplier write, or paid inference occurred. Pricing changes are not active on production.

## Image retention cleanup — 2026-10-11

At the user's request, retained the running images and the most recent rollback images without restarting any service. Retained main-app images are `sub2api1:deepseek-v4-pro-11f9760e-20261011` (running), `sub2api1:gateway-cors-8ff6cfb17-20261009` (previous deployed), and the same-day `sub2api1:account-rule-eb6b9fba1-20261011` candidate. Retained upstream-ops images are `upstream-ops:jws-browser-20261009` (running) and `upstream-ops:no-profit-20261009` (previous). `postgres:18-alpine` and `redis:8-alpine` remain retained for their running containers.

Removed only unreferenced historical tags: `sub2api1:deepseek-v32-cny-9c3c6fbd8-20261007`, `sub2api1:grok-build-fast-2d2326e2e-20261009`, `upstream-ops:profit-reconcile-20261008`, and `upstream-ops:sub2api-sync-home-20261008`. No `--force`, volume deletion or builder-cache pruning was used. The image inventory changed from 11 to 7 images; cleanup evidence and pre/post manifests are in `/opt/sub2api/backups/image-cleanup-20261010T180533Z`. All four containers stayed healthy with restart count zero; origin and upstream health checks returned 200.

## DeepSeek V4 Pro pricing and account-rule rollout — deployed 2026-10-11

The user authorized replacing the production app image after local verification and candidate construction. Production now uses `sub2api1:deepseek-v4-pro-11f9760e-20261011`, image ID `sha256:0d136bc6b2216825cf4431bcb3914b49e9ca21307697af6126626c647a770b73`, container `ad9ff96cc91660356d62007bf7c2ad9acaa81fbc4a4951b14f2019151582ebb4`, executable SHA256 `47bd98a37359060c96a7c66a41d5cf8044bb9a0ce395215398981e147b58e461`, and Sub2API `0.1.329` commit `11f9760e1896b5d9436ca2c1f18456f7a45cd955`. DeepSeek V4 Pro pricing is input `$4.50`, output `$13.50`, cache hit `$0.15` per million tokens; configured multipliers remain independent.

Before the app-only switch, the candidate applied the single pending additive migration `256_account_stats_rules_user_billing.sql`, adding `apply_to_user_billing BOOLEAN NOT NULL DEFAULT FALSE`. The migration completed successfully; no existing usage, balance, pricing-card or multiplier records were rewritten. Backup and runtime evidence are in `/opt/sub2api/backups/deepseek-v4-pro-rollout-20261010T175122Z`; the PostgreSQL custom-format dump passed `pg_restore -l` before the change.

Only the `sub2api` Compose image field changed and the app container was recreated with `--no-deps --pull never`. PostgreSQL and Redis container IDs/start times, data mounts, Caddy configuration, frontend pointer and persistent data were preserved. The new app is healthy with restart count zero; origin and public `/readyz` return HTTP 200, and startup logs contain no `ERROR`, `FATAL` or `PANIC` entries. Roll back the app only by restoring that backup's `docker-compose.before.yml` and recreating `sub2api`; do not restore the database for an image rollback. The additive migration remains applied during an app rollback.

## Gateway CORS — deployed 2026-10-09

User authorized local/GitHub/server synchronization and candidate build, explicitly withholding rollout authorization. API-key gateway routes now accept arbitrary browser origins without credentialed CORS; panel/session endpoints retain their configured restrictions. Missing/invalid API keys and cookie-only requests remain unauthorized. Source commit `8ff6cfb17f1aae0091c18b72672ff90606996c1d`; verification documentation commit `4725bf68e`; local main and GitHub `chenshangyi1/sub2api1/main` synchronized. See `docs/gateway-cors-candidate-20261009.md` for validation and limits.

Candidate tag `sub2api1:gateway-cors-8ff6cfb17-20261009`, image ID `sha256:5667a6fd1e60793eda3c830b6551ad9ddf37cd2e831c03133db166dc5f923139`, executable SHA256 `7fbf4c5058684f254b93d1e0e92f2c4c6506ea9a206648298ef37aa1203b4051`. Source staging remains `/opt/sub2api/build_context/billing-847ce5afa-full`; source backup `/opt/sub2api/backups/gateway-cors-source-20261009T043449Z`; evidence under `/opt/sub2api/build_context/outputs/gateway-cors-candidate`. Seventeen isolated HTTP checks and focused Go/race tests pass. The existing fresh-database migration naming failure for `027_usage_billing_consistency.sql` is outside this change; isolated validation against exported schema and migration ledger passes. No business records were copied or paid generation calls made.

The user subsequently approved deployment. Production now uses `sub2api1:gateway-cors-8ff6cfb17-20261009`, container `107ee164bb285c77e4392bf0d2c94f5a62e23323e5a07f6c969c54a452857f13`, started `2026-10-09T05:27:47.029166888Z`, healthy with restart count zero. The final switch at 13:27:46 Singapore time reached healthy at 13:27:53; 0.5-second origin sampling saw one unavailable sample. Two earlier attempts rolled back on strict checks (environment array order; one official reference price that also refreshed in the old executable after restart); total impact across attempts was not fully measured. Actual billing cards/group rates, PostgreSQL/Redis containers, Caddy/config and frontend pointer are preserved; only Compose app image changed. Backup `/opt/sub2api/backups/gateway-cors-rollout-20261009T052744Z` validates, and candidate read-only migration validation passed. Public CORS preflight 204, no-key/cookie-only 401, restricted paths 403 and origin/public readyz 200 passed. Roll back only by restoring that backup's `docker-compose.before.yml` to live Compose and recreating app-only with `--no-deps --pull never`; do not restore database for image rollback. See the synchronized rollout document for full evidence.

This handoff records the current operation, superseding old runtime assumptions only where explicitly verified below. No credentials are stored here.

- Local checkout: `/Users/chenmingxin/Downloads/sub2api-main 2`.
- User-selected GitHub destination: `https://github.com/chenshangyi1/sub2api1`, branch `main`; starting commit `847ce5afaa0d111e170406c32712fc7417d6f62d`.
- Authorized scope: correct SiliconFlow V3.2 base pricing, push source, and sync source to the server. Other DeepSeek models were audited, not repriced. Group/user multipliers are preserved.
- See `docs/deepseek-pricing-audit-20261007.md` for policy, verified sources, outstanding issues and test limitations.

## Verified server state before synchronization

- SSH origin: `root@216.22.13.148:22`, existing authentication via xinmc helper with host checking enabled.
- Running container ID: `0abe6b4db7624f9c7ec2e7fa4a8b587aefa97c59b530ba3c58f37ddc0cae6811`.
- Running image: `sub2api1:deepseek-standard-pricing-20261007-2`.
- StartedAt: `2026-10-06T19:53:49.955885423Z`; healthy, origin `/readyz` HTTP 200.
- Runtime Compose: `/opt/sub2api/deploy/docker-compose.local.yml`. Root checkout and Compose have intentional local changes; never reset or overwrite them.
- `/opt/sub2api/build_context` also has local modifications and an older Git baseline. Do not treat it as aligned with local main.
- Existing source staging tree `/opt/sub2api/build_context/billing-847ce5afa-full` matches the local starting billing source and test hashes. Reuse this tree for source-only synchronization; do not create another parallel source tree. It is not the canonical build root. The user subsequently authorized a candidate image build after local verification, but did not authorize switching the running service.

## Synchronization and rollback

Transfer only the committed changed source/tests and this operation's notes after checking the prior file hashes. Preserve all existing root, Compose, Caddy, environment, data and frontend paths. Save pre-change copies in a root-only `/opt/sub2api/backups/` directory and record the exact source commit in the staging tree's `SOURCE_REVISION`.

Recheck the container ID, image, start timestamp and readiness after transfer. They must remain unchanged. The source update does not activate the new prices; current production pricing continues until an explicitly authorized rollout.

For source rollback, use the recorded backup to restore the two previously existing changed files and remove only files newly introduced by this operation. Do not restore the database or change the running container. Before a future rollout, reconcile the canonical build root and review the supplier-specific pricing policy, migrations, known full-test failures and runtime rollback plan.

## Candidate build authorized after local verification

The candidate uses a locally cross-compiled linux/amd64, CGO-disabled binary with `embed timetzdata` tags and GOEXPERIMENT=jsonv2. Before compilation, all 2,883 backend Go/module files present in the running image were compared with local source: only this operation’s billing source and test differed, with the new V3.2 test present locally. All 221 embedded frontend files matched the running image byte for byte (excluding macOS metadata). No other runtime-source change is included.

Package the binary from `/opt/sub2api/build_context/outputs/deepseek-v32-candidate/sub2api` using `deploy/Dockerfile.deepseek-v32-candidate` from the canonical build root. Its Dockerfile-specific ignore file admits only that binary. The base must remain image ID `sha256:7ef469688d398f748e7cca40a8b105f3a6a4e2a506e0f4ec2953c83dbe6e79d7`, tag `sub2api1:deepseek-standard-pricing-20261007-2`. This preserves existing runtime dependencies and files while replacing only the executable. Record the candidate tag, source revision, binary hash and isolated version check after construction.

V3.2 regression coverage now includes the actual gateway calculation path with and without its pricing resolver, a real bill’s token counts, zero multiplier, zero-token/cache-only requests, partial/zero custom cards and concurrent group isolation. Expanded existing billing and tier tests must be reported separately from unrelated baseline failures; full-package test compilation still has the previously recorded blockers. Building a candidate does not activate prices. Do not restart the app, change Compose, mount production volumes into a candidate, or connect a candidate to production databases.

## Verified candidate result (built, not deployed)

- Candidate tag: `sub2api1:deepseek-v32-cny-9c3c6fbd8-20261007`.
- Image ID: `sha256:e8f9acf6dfcab55ace34151be7526456acdb0f2706592a3c9c4e4a44e9b18f9a`.
- Compiled source revision: `9c3c6fbd809441d657da26e9039f29ca8c3ab009`; application version 0.1.329; built at 2026-10-07T08:23:01Z. Later documentation-only commits do not change this executable.
- Executable SHA256: `b8cf200a71a28313d7978d8da0b3953756d5e14cbb9e925ce1bffa67e64d9f9b`. Local ELF and Go build information confirm linux/amd64, CGO_ENABLED=0 and embed/timetzdata. Git's dirty marker arises from the unrelated untracked api_key file; tracked source was clean at compilation.
- Candidate artifact/evidence directory: `/opt/sub2api/build_context/outputs/deepseek-v32-candidate`; `candidate_build.json`, `provenance.json`, `isolated-pricing-test.txt` and the executable are retained there.
- Source-sync backup for this verification/build step: `/opt/sub2api/backups/v32-candidate-20261007T082356Z`. Original price-patch backup remains `/opt/sub2api/backups/v32-source-20261007T081207Z`.
- Isolated candidate version execution and executable-hash check passed. Seven cross-compiled V3.2 regression tests passed in a read-only, network-disabled candidate container, with only its test executable mounted and no production configuration/data mounts.
- Runtime command, entrypoint, environment, healthcheck, working directory, user and exposed ports match the base image. Production container ID, image ID and StartedAt remain unchanged; healthy and origin /readyz returns HTTP 200. No Compose, database, group multiplier or static release change occurred.

This is a binary-overlay image: inherited `/app` Go source files are retained from the old image and are not the source of the new executable. For future builds, use the recorded GitHub revision/source staging, not the inherited `/app` source. The canonical server source tree's preexisting modifications were preserved; its candidate recipe consumes only the verified artifact.

The user later explicitly authorized rollout. On 2026-10-07T16:21:18Z, the app service was switched to `sub2api1:deepseek-v32-cny-9c3c6fbd8-20261007` after a read-only migration check passed. Backup and evidence are in `/opt/sub2api/backups/v32-rollout-20261007T161509Z`; the PostgreSQL custom-format dump was validated with the PostgreSQL container's `pg_restore -l`. Only the Compose app image field changed. The app became healthy at 2026-10-07T16:21:25Z; the measured readiness timeline had one unavailable sample during the single-container switch. The new container ID is `cefd4f8d52c1355bd0a3518fb80c5105fb8042c485a74cbab0d3acef6e99d423`, restart count 0, image ID `sha256:e8f9acf6dfcab55ace34151be7526456acdb0f2706592a3c9c4e4a44e9b18f9a`, and executable SHA256 `b8cf200a71a28313d7978d8da0b3953756d5e14cbb9e925ce1bffa67e64d9f9b`. Origin and public `/readyz` returned HTTP 200, public model-plaza reports V3.2 base prices 4/6/0.40 per million with group 16 rate 0.9, and other model prices, group rates, database, Redis, frontend symlink, mounts and runtime settings were preserved. Startup logs contained no ERROR/FATAL/PANIC lines. Rollback is to restore the backed-up Compose file and recreate only `sub2api`; do not restore the database unless a separate data-recovery event requires it.


## Grok build-fast candidate — 2026-10-09, built but NOT deployed

User authorized pushing `chenshangyi1/sub2api1/main`, server source sync and candidate image construction only. Source commit `2d2326e2ebc0f3cbb2d35dc54844a3cb120c46aa` sets only `grok-4.7-build-fast` base pricing to 4/12/1 per million, preserving explicit custom cards and multipliers. GitHub main was verified at this commit. Staging source `/opt/sub2api/build_context/billing-847ce5afa-full` matches 3,110 local Go/module/embedded-frontend files and has no extra Go source. Source backup: `/opt/sub2api/backups/grok-build-fast-source-20261009T074225Z`.

Candidate `sub2api1:grok-build-fast-2d2326e2e-20261009`, image ID `sha256:f434861458e9f97fd0338a85368337c2f07bbf770a890f3c30ed9364831b4ff2`, executable SHA256 `d0cd0a5ae90033c28321d94835554b5eebc910eec4dd83b7364c101a4caae247`. Compiled on the server with the inherited Go 1.27.0 toolchain, single-process/low-memory settings; candidate evidence, scripts and source manifest live in `/opt/sub2api/build_context/outputs/grok-build-fast-2d2326e2e`. Four focused pricing tests plus subcases passed in a read-only, network-disabled candidate container without production configuration/data mounts. Full suite limitations remain documented in `docs/grok-build-fast-local-pricing-20261009.md`.

Production container remains `107ee164bb285c77e4392bf0d2c94f5a62e23323e5a07f6c969c54a452857f13`, image `sub2api1:gateway-cors-8ff6cfb17-20261009`, StartedAt `2026-10-09T05:27:47.029166888Z`, healthy, restart 0, origin readyz 200. Compose/Caddy hashes and frontend pointer were unchanged; candidate runtime configuration matches its base. No deployment, migration, database or multiplier mutation occurred. Rollout needs explicit user approval and fresh migration/backup/runtime validation. Inherited image source is old; use the staging/Git revision for future builds.

SSH successfully bypassed Clash TUN using `ProxyCommand=/usr/bin/nc -4 -b en0 -G 10 %h %p`, preserving host checking and existing authentication. Verify the real physical interface again if networking changes; do not reuse a TUN-routed control socket for a bypass test. No global network/Clash/SSH changes were made.


## Historical image cleanup — 2026-10-09

User authorized removal of older images while retaining the current and previous versions. Removed 23 unused historical image tags without force and cleared unused Docker builder cache. Retained main-app current `gateway-cors-8ff6cfb17-20261009`, previous `deepseek-v32-cny-9c3c6fbd8-20261007`, and approved-for-build/unreleased `grok-build-fast-2d2326e2e-20261009`, plus upstream-ops current `profit-reconcile-20261008`, previous `sub2api-sync-home-20261008`, postgres:18-alpine and redis:8-alpine. Older rollback images named elsewhere in this document have been removed; do not use those tags as available local rollback targets. The immediate previous main-app rollback image remains available.

Evidence/inventory/deletion log: `/opt/sub2api/backups/image-cleanup-20261009T093740Z`; script `/opt/sub2api/build_context/image-cleanup-20261009.py`. Disk used decreased from 60,955,009,024 to 13,975,814,144 bytes (46,979,194,880 bytes released, about 43.75 GiB), with 87,240,400,896 bytes available and 14% usage. Docker reports 7 images and zero build cache. All four running containers retained identical IDs, image IDs, start timestamps, restart counts, health and mounts; Compose/Caddy hashes, frontend pointer and volume inventory were unchanged. Origin and public readyz returned 200. No restart, deployment, data/backup/source deletion or volume pruning occurred. Future builds may take longer because unused Docker build caches have been cleared; retained image layers remain usable.
