# Gemini 4 Argon candidate — 2026-10-14

## Scope and pricing

The user authorized local verification, GitHub synchronization, source staging synchronization and candidate image construction. Production deployment is explicitly excluded.

The accepted default card applies to `gemini-4-argon`, `gemini-4-argon-a`, `gemini-4-argon-instant`, their `models/` forms, and case/whitespace variants. Per million tokens: input USD 2, output USD 10, cache read USD 0.10, cache creation USD 0. Explicit channel/group cards and separately configured multipliers remain authoritative. The Argon base card takes precedence over the generic dynamic catalog.

Source: https://ai.google.dev/gemini-api/docs/pricing#gemini-4-argon, checked 2026-10-14. Google's Standard API table explicitly lists `gemini-4-argon`, implicit caching only, and the above promotional rates through 2026-11-30. It lists input/output/cache read USD 4/20/0.20 starting 2026-12-01. This candidate deliberately fixes the currently accepted card and does not schedule that increase. Review the default card before that date.

The `-a` and `-instant` names are the user's Antigravity supplier aliases, priced identically per the accepted screenshot/design; the official table does not independently establish separate prices for those names. Zero cache creation is the local card value for a model whose official API does not support explicit caching. This change covers ordinary token billing, not Google's Batch/Flex/Priority or grounding request fees.

## Source and verification

- Pricing source commit: `a969eeb314b3cae71b6646c900dbf9556fb55035`.
- Additional gateway/custom-card regression commit and compiled source: `94f8a21ef5b4364d93341fac30b337124f864e10`.
- GitHub destination: `chenshangyi1/sub2api1`, branch `main` (remote name `sub2api1`). The separate `origin` remote points to `chenshangyi1/sub2api`; it was not overwritten.
- Local focused Gemini run: 75 passing tests/subtests, command `GOEXPERIMENT=jsonv2 go test ./internal/service -run 'Test(GeminiArgon|BillingService_Gemini|PricingService_Gemini|DefaultPricingIncludesGemini)' -count=1 -v` from `backend`.
- Local Argon race run passed: `GOEXPERIMENT=jsonv2 go test -race ./internal/service -run 'Test(GeminiArgon|BillingService_Gemini4Argon)' -count=1`.
- Added coverage exercises the actual gateway token-cost function with and without its resolver, all six aliases plus a case/whitespace variant, zero/0.9/1.3 multipliers, a conflicting dynamic catalog, partial cards, explicitly free cards, custom cards and rejection of unknown suffixes by deterministic pricing identification. Tests do not record usage, change balances or call an inference provider.
- Linux/amd64 server and test executables compile with Go 1.26.6, CGO disabled and `GOEXPERIMENT=jsonv2`. The server also uses `embed timetzdata`, stripped debug information and version/revision/build-time ldflags. The embedded Go build metadata records the actual toolchain; do not infer it from older candidate recipes.
- Server source staging: `/opt/sub2api/build_context/billing-847ce5afa-full`; 3,428 backend source/module/migration/embedded-asset files matched the local manifest. No additional Go compilation inputs were found; macOS `._` metadata is excluded.
- Original Argon source backup: `/opt/sub2api/backups/argon-source-20261014T144227Z`; resumed synchronization backup: `/opt/sub2api/backups/argon-resume-20261014T153019Z`.

Full backend/frontend suites were not rerun for this focused backend addition. The existing full-suite failures recorded in `claude-jws-candidate-20261011.md` remain limitations; this operation does not claim the full suite is green.

## Candidate packaging

The compiled revision includes the previously committed, still undeployed Claude JWS billing changes. Compared with the running backend's source revision, this is not an Argon-only executable. Review the Claude JWS candidate document and this document together before any future rollout. There are no migration file differences from the running backend revision.

Package the verified binary over the pinned running image `sha256:0d136bc6b2216825cf4431bcb3914b49e9ca21307697af6126626c647a770b73`. The minimal Docker context contains only the recipe and executable. Build with command-local `DOCKER_BUILDKIT=0` and `--network=none`. The inherited Go source in the image belongs to the old base; future compilation must use the recorded Git revision or verified source staging.

- Verified tag: `sub2api1:gemini-4-argon-94f8a21ef-20261014`.
- Image ID: `sha256:c2da8fbaf6a37ceb1505277a5a3a3b7f5a150e752859af04752dc0a994229d38`.
- Binary SHA256: `fa1f3588cd69a9b842bacb93feb24e1cbbc1fcc94a9fa8ecb9c25c1d28216200`.
- Test executable SHA256: `bb2ca28af6fcf39208162d3f30801da20d547c1f146025185065c28e83333fa2`.
- Build timestamp: `2026-10-14T15:28:27Z`.
- Artifact/evidence destination: `/opt/sub2api/build_context/outputs/gemini-4-argon-candidate`.
- Scripts/provenance/source manifest: `/opt/sub2api/incoming/gemini-4-argon-94f8a21ef`.

## Verified result — built, NOT deployed

Image construction, embedded revision/version, executable hash and all 75 focused isolated tests/subtests passed. Containers used `--network=none`, a read-only filesystem, 0.5 CPU, 512 MiB and 128 PIDs; only the test executable and a public pricing fixture were mounted read-only. No production configuration, credentials, data volumes or databases were mounted/contacted.

The first isolated run reproduced a missing relative-path fixture in `TestDefaultPricingIncludesGemini36FlashRates`; Argon tests passed. Rerunning with the package's working directory and the matching local fixture fixed the environment. The fixture hash is `6f0bb924bb7509fd7ecbe3602d77c5a0d2f9d254810ad82ed6a8ee15eb51b43f`. The inherited runtime catalog and staging catalog instead have hash `d5bddaa053bc509653ca0b7cde4e3b57f1736795ba22f22e43db2cfdd98a3827`; these were preserved. Thus the full-file manifest above covers compilation inputs/migrations/embedded assets, not equality of every external runtime resource. The Argon default card is independent of that dynamic/fallback catalog. The unrelated Gemini fixture check validates the local fixture; it does not prove equality of the inherited runtime catalog. Failure evidence is retained as `isolated-billing-tests-missing-fixture.log`.

Candidate environment, command, entrypoint, working directory, user, healthcheck, exposed ports and declared volumes match the pinned base image. The final check confirmed all four production containers retain identical IDs, image IDs, start times, restart counts and health. Main app remains `sub2api1:deepseek-v4-pro-11f9760e-20261011`, container `ad9ff96cc91660356d62007bf7c2ad9acaa81fbc4a4951b14f2019151582ebb4`, healthy with restart count zero. Compose/Caddy hashes and frontend symlink are unchanged.

The initial build and successful final check observed 16 and 3 origin readiness samples respectively, all HTTP 200; final public `/readyz` also returned 200. No production rollout/restart, migration, billing data/card/multiplier write or paid inference occurred. New pricing is not yet active on production. Candidate construction does not trigger image cleanup.
