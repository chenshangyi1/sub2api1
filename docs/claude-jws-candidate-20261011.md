# Claude JWS candidate verification — 2026-10-11

Authorized scope: synchronize the reviewed local source to chenshangyi1/sub2api1/main and the existing server source staging tree, construct a candidate image, and verify it without switching production. No production restart, Compose change, database mutation, migration execution, or paid inference request is included.

Claude uses JWS upstream standard model prices. Fable 5/5.1 use 12/60 per million input/output, read 1.2/0.3; no additional 1.2 premium or Max 3x surcharge. Explicit custom cards remain authoritative. Group/user multipliers remain independently configured. Opus 5.5 actual JWS records verify input 4, output 20, write 5m 5, write 1h 8, read .2. Sonnet 5 and Haiku 4.5 have actual usage replay fixtures. Fable cache writes and some model/context/speed variants lack actual JWS evidence; see anthropic-billing-local-20261011.md.

## Local verification

- Selected production GoFiles plus billing tests: 496 test/subtest passes, 0 failures; two already recorded Grok baseline tests are explicitly skipped.
- Selected pricing, hold, protocol and replay paths with -race: 64 test/subtest passes.
- Handler balance preauthorization and API key route binding test files pass.
- Frontend model plaza table: 25 tests pass; main Vue typecheck passes; dev preview suite: 14 tests pass.
- No frontend source or embedded asset modification is included. Source-only ESLint has an existing parse failure in src/components/modelPlaza/PlazaHeroSection.vue:30; full lint also includes generated canvas assets (1031 errors, 472 warnings).
- Full frontend: 269 files / 1948 tests pass, 8 files / 12 tests fail and 1 unhandled rejection. Dev typecheck has 2 preexisting native_compaction_v2 fixture errors.
- Full backend unit invocation cannot be called green: service has duplicate declarations and stale account_stats_pricing tests, handler has removed model helpers/stale authorizer signatures, setup references removed helpers. Other observed failures are named below.
- Broader standalone model-plaza service test-file compilation also encounters missing/stale dependencies in existing gateway/scheduler test files; direct plaza reference/schedule checks in the financial regression and frontend presentation checks pass.
- golangci-lint is not installed on the local host. This verification does not claim complete lint or full-suite success.

## Observed full-suite failures

Backend:

- estConfigKeysAreEnvReachable
- estAllPlatformsIncludesEveryConcretePlatform
- estUpdateWithAccountBillingSettingsRollsBackWhenOutboxFails
- estAPIKeyRepositoryListByUserIDAttachesLastUsedIP
- estAPIKeyRepository_CreateWithLastUsedAt
- estAPIKeyRepository_UpdateLastUsed
- estAPIKeyRepository_UpdateLastUsedDeletedKey
- estAPIKeyRepository_UpdateLastUsedDBError
- estAPIKeyRepository_CreateDuplicateKey
- estAPIKeyRepository_GetByKeyForAuth_PreservesMessagesDispatchModelConfig_SQLite
- estLiveBalanceAdjustmentOutboxMigrationDefinesTransactionalSkip
- estAPIContracts
- estGatewayRoutesAdaptiveInboundDoesNot404LeafProtocols
- estAdaptiveHybridParentSQLStaysOutOfEmbeddedMigrations
- estAdaptiveUsageEvidenceMigration_EnforcesReservationAttemptUniqueness
- estAPIKeyAdaptivePassiveFailoverMigration

Frontend:

- src/api/__tests__/admin.accounts.stats.spec.ts > admin account stats API > uses an endpoint-specific 180s timeout
- src/i18n/__tests__/openaiFastPolicyLocales.spec.ts > OpenAI Fast/Flex policy locale keys > describes target and other-model actions without whitelist terminology
- src/components/layout/__tests__/AppHeader.spec.ts > AppHeader model plaza > keeps the model plaza entry visible on small screens
- src/components/payment/__tests__/PaymentProviderDialog.spec.ts > PaymentProviderDialog payment guide > rejects custom EasyPay method types with built-in payment prefixes
- src/views/admin/__tests__/GroupsView.codexManifest.spec.ts > GroupsView Codex manifest binding > preserves consecutive child updates on the reactive edit config
- src/views/admin/__tests__/SubscriptionsView.userUsageLink.spec.ts > admin subscription users > searches current users when assigning a subscription
- src/views/admin/__tests__/SubscriptionsView.userUsageLink.spec.ts > admin subscription users > clears the assignment user immediately when input changes to "another"
- src/views/admin/__tests__/SubscriptionsView.userUsageLink.spec.ts > admin subscription users > clears the assignment user immediately when input changes to ""
- src/views/admin/__tests__/SubscriptionsView.userUsageLink.spec.ts > admin subscription users > renders the user email as a link to that user filtered usage records
- src/views/admin/__tests__/SubscriptionsView.userUsageLink.spec.ts > admin subscription users > uses the user ID label for the usage link when username mode has no username
- src/views/user/__tests__/DashboardView.spec.ts > user DashboardView > 首屏不阻塞 refreshUser，最近用量只取 5 条
- src/views/admin/ops/components/__tests__/OpsErrorDetailModal.spec.ts > OpsErrorDetailModal > shows a WAF chip for Cloudflare 1010 bodies

## Image and synchronization plan

Use Go 1.27.0, linux/amd64, CGO_ENABLED=0, GOEXPERIMENT=jsonv2, embed/timetzdata; embed the existing locally verified assets. Build locally to avoid compiling on the live host. Package only the executable with deploy/Dockerfile.claude-jws-candidate, pinned to the current production image ID. Retain exact revision/hash/build information and cross-compiled financial test executable.

Before source writes, back up changed files in the existing /opt/sub2api/build_context/billing-847ce5afa-full tree, compare expected old hashes, and record SOURCE_REVISION. Preserve canonical /opt/sub2api checkout local changes. Candidate checks use read-only, network-disabled containers with limited CPU/memory and no production config/data/database mounts. Compare live container IDs, image IDs, start times/restarts, Compose/Caddy hashes and frontend symlink before/after. Final build evidence and hashes will be recorded after construction.
