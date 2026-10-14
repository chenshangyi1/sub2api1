# Gemini 4 Argon rollout — 2026-10-14

The user explicitly authorized switching to the verified candidate after reviewing its scope. This deploys the Argon default card and the previously undeployed Claude JWS billing changes described in the candidate documents.

## Deployed version

- Image: `sub2api1:gemini-4-argon-94f8a21ef-20261014`.
- Image ID: `sha256:c2da8fbaf6a37ceb1505277a5a3a3b7f5a150e752859af04752dc0a994229d38`.
- Container ID: `8c86dc27e51b692a6ddc3ea74dac51972d24a68f2a807ea09a7376f12d24399d`.
- Compiled source: `94f8a21ef5b4364d93341fac30b337124f864e10`; version 0.1.329.
- Executable SHA256: `fa1f3588cd69a9b842bacb93feb24e1cbbc1fcc94a9fa8ecb9c25c1d28216200`.
- Previous image retained for rollback: `sub2api1:deepseek-v4-pro-11f9760e-20261011`, ID `sha256:0d136bc6b2216825cf4431bcb3914b49e9ca21307697af6126626c647a770b73`.

## Backup and preflight

Backup/evidence directory: `/opt/sub2api/backups/argon-rollout-20261014T155011Z` (root-only). It contains the 2,337,168-byte PostgreSQL custom-format dump, validated with PostgreSQL 18 `pg_restore -l`, the original Compose/Caddy/application configuration, container/runtime snapshots, financial configuration snapshots and operation logs. Archive validation checks readability/format, not an isolated full restore.

Candidate `--check-migrations` passed against the existing database in a read-only container with production app mounts read-only. That code path only reads the migration ledger and avoids Redis, workers and startup seeders. No pending migrations were found and no migration command was run in write mode. Candidate image ID, executable checksum and embedded source version were independently verified before switching.

## Switch and actual interruption

Only the `sub2api` image field in live Compose was changed; normalized Compose configuration otherwise matched the backup. The command was app-only `docker compose -p deploy -f /opt/sub2api/deploy/docker-compose.local.yml up -d --no-deps --pull never sub2api`.

- Switch started: **2026-10-14 23:51:29 Singapore time** (`15:51:29.046485Z`).
- Healthy confirmed: **2026-10-14 23:57:37 Singapore time** (`15:57:37.084173Z`).
- Time from switch start to healthy: **368.038 seconds**.
- Origin sampling interval: 0.5 seconds plus request time; 709 unavailable samples, with **362.886 seconds** from first unavailable sample to the first successful sample after the last failure.
- One separate public sample during the switch returned 502. The origin sampling span is not an exact measurement of every public request's outage duration.

The prior 5–15 second interruption estimate was incorrect: it omitted the already configured `SERVER_SHUTDOWN_TIMEOUT_SECONDS=600`, `SUB2API_UPDATE_DRAIN_TIMEOUT_SECONDS=610` and Compose `stop_grace_period: 610s`. The old app closed its listener while waiting for in-flight requests; logs showed a disconnected streaming request still being drained for billing. Its actual exit preceded the configured maximum. No manual kill or drain-timeout reduction was performed. These settings remain intact; future rollout estimates must account for them, or use a separately planned deployment strategy that keeps a serving instance available.

## Post-deployment checks

- New image, executable checksum and embedded revision/version match the verified candidate.
- Healthy with restart count 0; origin/public `/readyz` and origin `/api/v1/model-plaza` returned HTTP 200.
- 34 Argon-specific test/subtest events passed in a read-only, network-disabled container using the deployed image and only the test executable mounted read-only. Earlier candidate/local focused Gemini checks each passed 75 events. No paid inference probe was made.
- Startup logs contained zero ERROR/FATAL/PANIC lines in the inspected interval.
- PostgreSQL, Redis and upstream-ops retained their container IDs, images, start times and restart counts; all remain healthy.
- App environment, command, entrypoint, user, working directory, healthcheck, declared volumes, mounts and published port bindings were preserved. Env ordering was normalized for comparison.
- Caddy/application configuration and the frontend release pointer were unchanged.
- Read-only before/after comparisons confirmed group rate/default model cards, user-group rate overrides and channel cards were unchanged. User traffic and ordinary background usage processing continued; balances/usage were not frozen or rewritten for this rollout.

## Image retention and rollback

After successful checks, removed only these unreferenced, superseded main-app tags without force: `sub2api1:gateway-cors-8ff6cfb17-20261009`, `sub2api1:claude-jws-2b4f5d51a-20261011`, `sub2api1:account-rule-eb6b9fba1-20261011`. The current Argon image and exact pre-switch DeepSeek image remain available. Other services' images, volumes, build caches, data, backups and source were preserved. Runtime/config fingerprints stayed identical across cleanup; origin/public readiness both returned 200. Available disk after cleanup: 81,656,016,896 bytes. Manifests and cleanup plan/results are in the backup directory.

For an app rollback, restore that directory's `docker-compose.before.yml` to live Compose and run the same app-only `up -d --no-deps --pull never sub2api`. Recheck health, executable identity and public readiness. The existing drain settings can again lengthen a rollback. **Do not restore the database for an image rollback.**
