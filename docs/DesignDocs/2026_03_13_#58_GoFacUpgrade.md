# GoFac Go 1.22 → 1.26 Upgrade Plan

## Metadata

- Date: 2026-03-13
- Owner: TaBSRest
- Superseding: N/A
- Preceding: N/A
- Status: Draft

## Brief & Problem Statements

GoLang 1.24 has reached end of support. GoFac currently targets Go 1.22.0 per `go.mod`. To stay on a supported toolchain and unblock downstream consumers (MobileBFF, BE, WebHook), the upgrade path is:

1.22 (tag `1.22.0.8.3`) → 1.23 (tag `1.23.0.8.3`) → 1.24 (tag `1.24.0.8.3`) → 1.25 (tag `1.25.0.8.3`) → 1.26 (tag `1.26.0.8.3`)

Each minor version bump is its own PR on `main` (trunk-based development). Downstream consumers are updated in separate PRs after `1.26.0.8.3` is tagged.

## Considerations and Constraints

- Each version bump touches at most 3 files: `go.mod`, `go.sum`, `.github/workflows/go.yml`.
- All existing tests must pass before each merge.
- The CI workflow (`.github/workflows/go.yml`) does not currently pin a Go version; each PR introduces or updates the `actions/setup-go` pin.
- GoFac has no CGo, no Wasm, no `//go:linkname`, no timers, no crypto usage — the vast majority of breaking changes in each release do not apply.
- Each version gets its own tag (`1.22.0.8.3` through `1.26.0.8.3`) so downstream consumers can pin any intermediate version if needed.

## Approach to Problems

For each version: review breaking changes, confirm no GoFac code changes are needed (or make them), bump `go.mod`, run `go mod tidy`, update the CI Go pin, run the full test suite, open PR, merge, tag if required.

## Solution

### Architectural Overview

GoFac (`github.com/TaBSRest/GoFac`) is a dependency injection container library with two direct dependencies: `github.com/google/uuid v1.6.0` and `github.com/stretchr/testify v1.9.0`. No CGo, no Wasm, no generated code. The upgrade is entirely a toolchain and module file change.

### Design Efficacy

Stepping through each minor version individually:

- Isolates which version introduces a regression if tests start failing.
- Keeps each PR minimal and reviewable.
- Gives downstream consumers a tagged pin at every Go minor version (`1.22.0.8.3` → `1.26.0.8.3`), so they can adopt at their own pace.

### Limitations

- Green Tea GC (enabled by default in Go 1.26) is a binary-level change; downstream callers control it. The library has no influence over it.
- This plan does not cover updating downstream repos' full test suites — only their GoFac dependency pin.

### Edge Cases

- `go mod tidy` after a version bump may add or remove a `toolchain` line in `go.mod` (expected per Go 1.25's new toolchain line behavior). Commit the diff as-is.
- If the CI runner's default Go lags behind the pinned version, the `actions/setup-go` step handles it.

## Implementation Steps

Each step below is one PR. Branch naming: `feat/{codex}/#{ticket}-{description}`. PR title: `feat(GoFac): #{ticket} Upgrade Go to <version>`.

---

### Pre-step — Tag current state as `1.22.0.8.3`

No PR needed — tag `main` at the current commit before any changes.

- [ ] Tag `1.22.0.8.3` on the current `main` HEAD.

---

### PR 1 — Upgrade to Go 1.23 + tag `1.23.0.8.3`

**Files changed:** `go.mod`, `go.sum`, `.github/workflows/go.yml` (3 files)

- [ ] Research Go 1.23 breaking changes (see Appendix). No GoFac code changes required.
- [ ] Update `go.mod`: `go 1.22.0` → `go 1.23`.
- [ ] Run `go mod tidy`; commit `go.mod` and `go.sum`.
- [ ] Add `actions/setup-go` step to `.github/workflows/go.yml` with `go-version: "1.23"`.
- [ ] Run `make checkCoverage` locally and confirm all tests pass.
- [ ] Open PR → merge to `main`.
- [ ] Tag `1.23.0.8.3` on the merge commit.

---

### PR 2 — Upgrade to Go 1.24 + tag `1.24.0.8.3`

**Files changed:** `go.mod`, `go.sum`, `.github/workflows/go.yml` (3 files)

- [ ] Research Go 1.24 breaking changes (see Appendix). No GoFac code changes required.
- [ ] Update `go.mod`: `go 1.23` → `go 1.24`.
- [ ] Run `go mod tidy`; commit `go.mod` and `go.sum`.
- [ ] Update `.github/workflows/go.yml` Go pin to `"1.24"`.
- [ ] Run `make checkCoverage` locally and confirm all tests pass.
- [ ] Open PR → merge to `main`.
- [ ] Tag `1.24.0.8.3` on the merge commit.

---

### PR 3 — Upgrade to Go 1.25 + tag `1.25.0.8.3`

**Files changed:** `go.mod`, `go.sum`, `.github/workflows/go.yml` (3 files)

- [ ] Research Go 1.25 breaking changes (see Appendix).
- [ ] Audit error-handling paths for the nil pointer check fix: confirm no code in GoFac dereferences a potentially-nil pointer before an adjacent error check.
- [ ] Update `go.mod`: `go 1.24` → `go 1.25`.
- [ ] Run `go mod tidy`; commit `go.mod` and `go.sum` (the `toolchain` line behavior changes in 1.25 — commit whatever `tidy` produces).
- [ ] Update `.github/workflows/go.yml` Go pin to `"1.25"`.
- [ ] Run `make checkCoverage` locally and confirm all tests pass.
- [ ] Open PR → merge to `main`.
- [ ] Tag `1.25.0.8.3` on the merge commit.

---

### PR 4 — Upgrade to Go 1.26 + tag `1.26.0.8.3`

**Files changed:** `go.mod`, `go.sum`, `.github/workflows/go.yml` (3 files)

- [ ] Research Go 1.26 breaking changes (see Appendix). No GoFac code changes required.
- [ ] Run `go fix ./...` locally and review any suggested modernizer rewrites before editing.
- [ ] Update `go.mod`: `go 1.25` → `go 1.26`.
- [ ] Run `go mod tidy`; commit `go.mod` and `go.sum`.
- [ ] Update `.github/workflows/go.yml` Go pin to `"1.26"`.
- [ ] Run `make checkCoverage` locally and confirm all tests pass.
- [ ] Open PR → merge to `main`.
- [ ] Tag `1.26.0.8.3` on the merge commit.

---

### PR 5 — Update MobileBFF to GoFac `1.26.0.8.3`

**Files changed:** `go.mod`, `go.sum` in the MobileBFF repo (2 files)

- [ ] Update `require github.com/TaBSRest/GoFac` to `v1.26.0.8.3`.
- [ ] Run `go mod tidy`.
- [ ] Confirm the downstream test suite passes.
- [ ] Open PR → merge.

---

### PR 6 — Update BE to GoFac `1.26.0.8.3`

**Files changed:** `go.mod`, `go.sum` in the BE repo (2 files)

- [ ] Update `require github.com/TaBSRest/GoFac` to `v1.26.0.8.3`.
- [ ] Run `go mod tidy`.
- [ ] Confirm the downstream test suite passes.
- [ ] Open PR → merge.

---

### PR 7 — Update WebHook to GoFac `1.26.0.8.3`

**Files changed:** `go.mod`, `go.sum` in the WebHook repo (2 files)

- [ ] Update `require github.com/TaBSRest/GoFac` to `v1.26.0.8.3`.
- [ ] Run `go mod tidy`.
- [ ] Confirm the downstream test suite passes.
- [ ] Open PR → merge.

---

## Appendix

### Breaking Change Research Notes

#### Go 1.23 Breaking Changes

**Source:** <https://go.dev/doc/go1.23>

| Change | Details | GoFac Impact |
|--------|---------|--------------|
| `time.Timer`/`time.Ticker` channel capacity | Changed from 1 to 0 for modules on `go 1.23+`; polling `len(ch)` breaks. Revert: `GODEBUG=asynctimerchan=1`. | **None** — GoFac uses no timers. |
| `//go:linkname` restriction | Cannot access internal stdlib symbols without explicit annotation; `-checklinkname=0` disables. | **None** — no `//go:linkname` in GoFac. |
| `crypto/tls` — 3DES removed from defaults | Revert: `GODEBUG=tls3des=1`. | **None** — no TLS usage. |
| `net/http` `ServeContent` header stripping | Error responses strip cache/encoding headers. | **None** — no HTTP serving. |
| `crypto/x509` — `CreateCertificateRequest` signature verification | Now returns error on bad signature. | **None** — no X.509 usage. |
| New `godebug` directive in `go.mod` | Opt-in; pins GODEBUG settings per module. | **None required.** |
| `reflect.DeepEqual` — `netip.Addr` | IPv4 vs. IPv4-mapped addresses now distinguished. | **None** — no `netip` usage. |
| Runtime traceback indentation | Second+ panic lines now tab-indented. | **None** — no tests assert on panic traceback text. |

**Verdict:** No code changes required.

---

#### Go 1.24 Breaking Changes

**Source:** <https://go.dev/doc/go1.24>

| Change | Details | GoFac Impact |
|--------|---------|--------------|
| `crypto/rand.Read` panics instead of returning error | Previously returned error on failure; now panics. | **None** — no `crypto/rand` usage. |
| `crypto/rsa` — minimum 1024-bit key | Keys < 1024 bits rejected. | **None** — no RSA usage. |
| `crypto/x509` — SHA-1 removed | `x509sha1` GODEBUG removed. | **None** — no X.509 usage. |
| `math/rand.Seed` is a no-op | Top-level `Seed()` silently ignored. Revert: `GODEBUG=randseednop=0`. | **None** — no `math/rand` seeding. |
| `runtime.SetFinalizer` deprecated | `runtime.SetFinalizer` is deprecated in Go 1.24. `runtime.AddCleanup` (added in Go 1.23) is the recommended alternative. | **Low** — `pkg/Container/Container.go` uses `SetFinalizer`; still works, migration to `AddCleanup` is optional. |
| `sync.Map` — Swiss Tables / hash trie internals | Performance improvement; public API unchanged. | **None** — GoFac uses `sync.Map` via public API only. |
| New `tool` directive in `go.mod` | Replaces blank-import `tools.go` pattern. | **None** — GoFac has no `tools.go`. |
| `go test -json` includes build events | New `Action: "build"` entries appear in JSON output. | **Low** — verify no CI step parses raw `go test -json` output. |
| Linux kernel minimum raised to 3.2 | Affects CI runners (all modern runners satisfy this). | **None** — verify runner OS. |

**Verdict:** No code changes required. Optionally migrate `SetFinalizer` → `AddCleanup` in a future PR.

---

#### Go 1.25 Breaking Changes

**Source:** <https://go.dev/doc/go1.25>, <https://go.dev/blog/go1.25>

| Change | Details | GoFac Impact |
|--------|---------|--------------|
| Nil pointer check fix | Compiler enforces spec-correct nil checks; previously the check could be delayed past an error check. | **Low** — audit error-handling paths before bumping. |
| GOMAXPROCS default (Linux cgroup) | Defaults to cgroup CPU limit; dynamically adjusts at runtime. | **None** — library does not set GOMAXPROCS. |
| `encoding/json` error text | Internal reimplementation; error wording may differ. | **None** — no `encoding/json` usage in GoFac. |
| `crypto/elliptic` undocumented methods removed | `Inverse`, `CombinedMult` removed. | **None** — no crypto/elliptic usage. |
| macOS minimum version raised to macOS 12 | Development/CI constraint only. | **None**. |
| `toolchain` line no longer auto-added | `go mod tidy` no longer appends `toolchain` pin. | **Low** — expect `go.mod` diff; commit it. |
| `testing/synctest` graduates to GA | Old `GOEXPERIMENT=synctest` API still present but removed in 1.26. | **None** — GoFac does not use synctest. |
| Generic "core types" removed from spec | Spec cleanup; behavior unchanged, error messages improve. | **None**. |

**Verdict:** No code changes required. Audit nil pointer paths; commit expected `go.mod` diff.

---

#### Go 1.26 Breaking Changes

**Source:** <https://go.dev/doc/go1.26>, <https://go.dev/blog/go1.26>

| Change | Details | GoFac Impact |
|--------|---------|--------------|
| Green Tea GC enabled by default | 10–40% GC overhead reduction. Opt-out: `GOEXPERIMENT=nogreenteagc` (removed in 1.27). | **None** — binary-level change; downstream callers control it. |
| `testing/synctest` old experimental API removed | `GOEXPERIMENT=synctest` variant from Go 1.24 gone. | **None** — GoFac does not use synctest. |
| `windows/arm` (32-bit) port removed | `GOOS=windows GOARCH=arm` unsupported. | **None**. |
| Bootstrap requirement raised to Go 1.24.6 | Building the Go toolchain itself requires 1.24.6+. | **None** — toolchain builder concern only. |
| `go fix` modernizer passes | May suggest or auto-apply safe rewrites. | **Low** — run `go fix ./...`, review output, commit if desired. |
| Self-referential generic constraints | Additive language feature. | **None**. |
| ELFv1 on ppc64 deprecated | Last release with ELFv1; ELFv2 in 1.27. | **None**. |

**Verdict:** No code changes required. Run `go fix ./...` and review.

---

### Diagrams

N/A — toolchain upgrade with no architectural changes.

### Other Solutions

#### Batch the Go version bumps into one PR

**Advantages:** Fewer PRs.

**Disadvantages:** Violates the ≤ 10 file / ≤ 300 line trunk-based rule; harder to bisect if a regression appears.

**Reasons for Rejection:** The organization rule states one logical work unit per PR.

#### Skip 1.23 and 1.24, jump from 1.22 directly to 1.25

**Advantages:** Fewer PRs; faster path to the tagged releases.

**Disadvantages:** Harder to isolate which version introduced a regression. Skips the incremental safety net.

**Reasons for Rejection:** Stepping through each version is safer and aligns with trunk-based development principles.
