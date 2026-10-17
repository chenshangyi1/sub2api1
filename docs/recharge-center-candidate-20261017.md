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
