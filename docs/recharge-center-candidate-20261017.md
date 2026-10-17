# Recharge center candidate - 2026-10-17

## Scope

This change adds an authenticated `/recharge` page that embeds the configured card shop at `https://wzyp.cn/shop/6B2BP3Q3`. Users buy a card in the shop, then use the existing `/redeem` page. The iframe receives only `ui_mode=embedded`, the current theme, and the locale. No account identifier, session token, source URL, purchase, payment callback, or balance write is sent to the shop.

The backend CSP allowlist includes `https://wzyp.cn` for `frame-src`. The active production Caddy policy still needs the same one-domain addition before a future frontend publication; no production Caddy file was installed or reloaded in this operation.

## Verification

- Focused frontend suite: 3 files, 52 tests passed.
- Frontend typecheck: passed.
- Frontend production build: passed; output is `backend/internal/web/dist`.
- Middleware CSP red/green: test failed without the new allowlist entry, then passed with `GOEXPERIMENT=jsonv2`.
- Local console preview at `http://127.0.0.1:4327/recharge`: authenticated demo fixture rendered the new navigation, controls, and the live shop iframe. The iframe rendered the merchant shop and six card products without a purchase.
- The local preview was stopped after verification. Temporary preview config files were removed.

## Rollout boundary

This is source synchronization and an isolated candidate build only. The running `sub2api` container, database, Redis, upstream-ops container, frontend symlink, and Caddy configuration remain unchanged. Do not point the production frontend symlink or reload Caddy until the Caddy `frame-src` addition is reviewed and a separate rollout is authorized.

## Candidate and synchronization results

- Runtime source commit: `f7077c3f05daa28bab36445022a6efe19df3e28e`, pushed to `chenshangyi1/sub2api1` on `main` and `codex/recharge-center`. Later documentation commits do not change the executable.
- Server artifacts: `/opt/sub2api/incoming/recharge-f7077c3f0-20261017`; this directory contains the source/frontend archives, manifests, build logs, isolated checks and candidate Caddy configuration.
- Existing source staging: `/opt/sub2api/build_context/billing-847ce5afa-full`. Twelve changed files were checked against their parent-commit hashes before copying; unrelated files were preserved. Backup: `/opt/sub2api/backups/recharge-source-20261017T083230Z`; `source-sync-manifest.json` records before/after SHA256 values. `SOURCE_REVISION` records the runtime source commit.
- Candidate image: `sub2api1:recharge-center-f7077c3f0-20261017`, ID `sha256:02aea3438eb6fb6b0e7b8368a1f54155b2fa3d08fb05ec0aa99e56f2d17b1d53`. Executable SHA256: `a6b7f8c616de844a98d6de6f920958a6ce87f212bab3b6b67a1a76f9d024d37c`.
- Packaging replaces the executable over production base `sha256:c2da8fbaf6a37ceb1505277a5a3a3b7f5a150e752859af04752dc0a994229d38`. Inherited image source is historical; rebuild using the recorded Git revision. Runtime entrypoint, command, environment, user, working directory, healthcheck and exposed ports match the base.
- Isolated checks: CSP middleware tests, version/revision and executable hash passed. A network-disabled, read-only, resource-limited container served `/recharge` with the shop CSP and `/setup/status` with `needs_setup: true`. No production configuration, database, Redis or data volume was attached. This checks empty-data startup, not a full production database integration or real payment flow.
- Candidate Caddy file: `Caddyfile.candidate`, SHA256 `931009e52c9b47b4e836963fe1fafdde6d26835128a7b4803b45a9fc6b5dd091`; Caddy validation passed. The active configuration was not replaced or reloaded.
- Production baseline preserved: app container `8c86dc27e51b692a6ddc3ea74dac51972d24a68f2a807ea09a7376f12d24399d`; Caddy SHA256 `1209df7f671c9bde3d966939b78657a96ac3daf1babed8ac434e7dff9cd6eb99`; Compose SHA256 `e485f37a6728a1cbe505b92ea099bae8b6bcc165c2aa1e8123859cd84f4f79ff`; frontend target `/opt/sub2api/frontend-releases/20261006T014419Z-profit-ui`. All four production containers retained IDs, images, start times, restart counts and healthy status; origin/public readiness returned HTTP 200.

No production image switch, migration, registration-balance change, purchase, payment, redemption or image cleanup was performed. Future frontend publication must install the built assets and the validated Caddy CSP addition together; merely building this image does not publish the externally served frontend.
