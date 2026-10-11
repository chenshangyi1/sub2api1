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

- TestConfigKeysAreEnvReachable
- TestAllPlatformsIncludesEveryConcretePlatform
- TestUpdateWithAccountBillingSettingsRollsBackWhenOutboxFails
- TestAPIKeyRepositoryListByUserIDAttachesLastUsedIP
- TestAPIKeyRepository_CreateWithLastUsedAt
- TestAPIKeyRepository_UpdateLastUsed
- TestAPIKeyRepository_UpdateLastUsedDeletedKey
- TestAPIKeyRepository_UpdateLastUsedDBError
- TestAPIKeyRepository_CreateDuplicateKey
- TestAPIKeyRepository_GetByKeyForAuth_PreservesMessagesDispatchModelConfig_SQLite
- TestLiveBalanceAdjustmentOutboxMigrationDefinesTransactionalSkip
- TestAPIContracts
- TestGatewayRoutesAdaptiveInboundDoesNot404LeafProtocols
- TestAdaptiveHybridParentSQLStaysOutOfEmbeddedMigrations
- TestAdaptiveUsageEvidenceMigration_EnforcesReservationAttemptUniqueness
- TestAPIKeyAdaptivePassiveFailoverMigration

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

Before source writes, back up changed files in the existing /opt/sub2api/build_context/billing-847ce5afa-full tree, compare expected old hashes, and record SOURCE_REVISION. Preserve canonical /opt/sub2api checkout local changes. Candidate checks use read-only, network-disabled containers with limited CPU/memory and no production config/data/database mounts. Compare live container IDs, image IDs, start times/restarts, Compose/Caddy hashes and frontend symlink before/after. Verified build results and hashes are recorded below.

## Verified candidate result — built, NOT deployed

- Candidate tag: `sub2api1:claude-jws-2b4f5d51a-20261011`.
- Image ID: `sha256:017b5b91dd650e44e9cc395cd5edfeb91e4e9eab9d13696e651daec0edc519c1`; base image ID: `sha256:0d136bc6b2216825cf4431bcb3914b49e9ca21307697af6126626c647a770b73`.
- Compiled source: `2b4f5d51a926256f2c7727b42f98b2e301f0e273`, application 0.1.329, built at 2026-10-11T13:49:46Z. Final documentation/packaging-only commits do not change the compiled source or executable.
- Executable SHA256: `f809796957a4fb1b089efb56ec3fae949656983b69a57c0e06d6213ebba80717`.
- Server artifacts/evidence: `/opt/sub2api/build_context/outputs/claude-jws-candidate`.
- Source staging: `/opt/sub2api/build_context/billing-847ce5afa-full`; source backup: `/opt/sub2api/backups/claude-jws-source-20261011T135045Z`.

Isolated candidate version/hash checks and 64 selected billing test/subtest events passed in read-only, network-disabled containers limited to 0.5 CPU, 512 MiB and 128 PIDs. Only the test executable was mounted read-only; production config, credentials, data and databases were not mounted or contacted. Candidate runtime environment, command, entrypoint, working directory, user, healthcheck, ports and declared volumes match the base image.

All four live containers retain the exact pre-operation IDs, image IDs, StartedAt, restart counts and health status. Main app remains `sub2api1:deepseek-v4-pro-11f9760e-20261011`, container `ad9ff96cc91660356d62007bf7c2ad9acaa81fbc4a4951b14f2019151582ebb4`, started `2026-10-10T17:59:13.97321849Z`, healthy, restart 0. Compose/Caddy hashes and frontend pointer are unchanged. Nine origin readiness samples during successful build/verification returned 200; final origin and public readiness both returned 200. Sampling is evidence of observed availability, not proof about every intervening request. No runtime switch, restart, migration, billing-data or multiplier write occurred.

Before packaging, 3,427 local backend source/module/migration/embedded-asset hashes matched staging with no extra source files. Canonical `/opt/sub2api` local changes were preserved. This image replaces only `/app/sub2api`; inherited Go files remain from its base and must not be used as the new executable's source. Future compilation must use the recorded Git revision/staging. The unrelated local `api_key` was excluded from commits and transfers; its presence accounts for Go's vcs.modified marker while tracked backend source was clean at compilation.

Packaging uses an explicit minimal context containing the recipe and executable at `outputs/claude-jws-candidate/sub2api`, mode 0755. Run the existing `build-candidate.py` against that context with command-local `DOCKER_BUILDKIT=0` and `--network=none`; the server legacy builder accepts the pinned local image ID and preserves COPY file permissions. Do not use the entire server source/backup tree as build context. The committed Dockerfile-specific ignore file is for builders supporting it; legacy packaging relies on the minimal context instead. No global Docker configuration was changed.

Full-suite failures listed above remain unresolved. This candidate is prepared and verified within the stated limits; it is not a globally green release and no application startup against a production database was attempted. New pricing is not active on the server.
