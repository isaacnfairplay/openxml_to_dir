# AGENTS.md

*Project: ooxmlx — CGO-free OpenXML Extractor and XML Prettifier (Go 1.23+)*

---

## 🧭 Purpose

This document defines the **agents**, **responsibilities**, and **standards** for the `ooxmlx` project.
It ensures consistent development, testing, release, and maintenance of a **pure standard library Go** command-line tool that securely extracts and prettifies XML documents from `.pptx`, `.docx`, and `.xlsx` files.

---

## 🧩 Architecture Agents

### 1. Core Extractor Agent (`internal/extract`)

**Responsibility:**
Implements the orchestration logic for reading the zip archive, delegating to path resolvers, transformers, and formatters.
Adheres to **SOLID principles**:

* **Single Responsibility:** Orchestrates extraction only.
* **Open/Closed:** Extend behavior via new transformers.
* **Liskov Substitution:** Any `Transformer` implementation can replace another.
* **Interface Segregation:** Works only with required interfaces (`Formatter`, `Resolver`, `Transformer`).
* **Dependency Inversion:** Depends on abstractions, not implementations.

**Guarantees:**

* Zip entries are extracted exactly once.
* XML detection and formatting are side-effect-free.
* Path safety validated before writing.

---

### 2. Path Safety Agent (`internal/pathsafe`)

**Responsibility:**
Prevents Zip Slip attacks by enforcing all writes remain inside the destination directory.

**Guarantees:**

* Rejects any extraction path resolving outside base.
* Uses absolute path comparison for platform-agnostic safety.
* Tested for traversal edge cases (`../../evil.txt`).

---

### 3. XML Formatting Agent (`internal/xmlutil`)

**Responsibility:**
Handles XML parsing, indentation, and serialization using only `encoding/xml`.

**Capabilities:**

* Cheap prefilter detects XML-like content.
* Parses token stream deterministically.
* Pretty-prints XML while preserving semantic correctness.
* Supports configurable indentation, encoding, and declaration control.

**Testing Coverage:**

* Unit tests for valid and invalid XML paths.
* Fuzzing ensures robustness under random input.
* Benchmarks validate performance under small/large documents.

---

### 4. Transformation Agent (`internal/transform`)

**Responsibility:**
Defines the `Transformer` interface and built-in transformations (e.g., NBSP normalization).

**Composition:**

* `Composite`: Chains multiple transformers in order.
* `ReplaceNBSP`: Rewrites U+00A0 to spaces.

**Guarantees:**

* Stateless, deterministic transformations.
* Easily extensible (Open/Closed).

---

### 5. ZIP Reader Agent (`internal/zipwrap`)

**Responsibility:**
Wraps `archive/zip` to ensure a single open handle per archive (no repeated reopens).

**Guarantees:**

* Minimal I/O overhead.
* Clean closure of file descriptors.
* Provides a simple iterable list of `zip.File` entries.

---

### 6. CLI Agent (`cmd/ooxmlx`)

**Responsibility:**
Implements the end-user interface and command-line parser.

**Flags:**

```
-o string          Destination directory (required)
-overwrite         Allow extraction into non-empty dest
-indent string     Indentation (default "  ")
-encoding string   XML encoding (default "utf-8")
-fix-nbsp          Normalize non-breaking spaces in XML
-quiet             Suppress logs
-version           Print version info
```

**Guarantees:**

* Clean exit codes.
* Prints version metadata from `internal/buildinfo`.
* Structured output logs for automation.

---

### 7. Build Metadata Agent (`internal/buildinfo`)

**Responsibility:**
Holds build-time metadata injected with `-ldflags` during releases.

**Injected Variables:**

* `Version` (from git tag)
* `Commit`  (from SHA)
* `Date`    (UTC build date)

**Printed via:**
`ooxmlx --version`

---

## ⚙️ CI/CD Agents

### Continuous Integration (CI)

**Workflow:** `.github/workflows/ci.yml`

**Responsibilities:**

* Run `go vet`, `go test -race -cover`.
* Publish coverage artifacts.
* Enforce no CGO (`CGO_ENABLED=0`).
* Execute fuzz targets for 20s to catch panics.

### Continuous Delivery (Release)

**Workflow:** `.github/workflows/release.yml`

**Responsibilities:**

* Triggered on `v*.*.*` tags.
* Builds for Linux, macOS, Windows (amd64/arm64).
* Injects build metadata via `-ldflags`.
* Zips artifacts and computes SHA256SUMS.
* Publishes GitHub Release with binaries and checksums.

**Guarantees:**

* Deterministic builds (no external deps).
* Cross-platform reproducibility.

---

## 🧪 Testing Agents

| Test Type       | Target                    | Method                 | Purpose                         |
| --------------- | ------------------------- | ---------------------- | ------------------------------- |
| Unit            | `pathsafe`, `xmlutil`     | `go test -race`        | Validate core logic correctness |
| Integration     | `extract`                 | in-memory zip archives | Full end-to-end coverage        |
| Fuzz            | `xmlutil.FuzzTryPrettify` | Go fuzzing             | Detect parser or panic issues   |
| Benchmarks      | `xmlutil`                 | `go test -bench`       | Performance measurement         |
| Static Analysis | CI                        | `go vet`               | Lint and verify idioms          |

---

## 📦 Release Agent

**Trigger:** `git tag vX.Y.Z && git push origin vX.Y.Z`
**Outputs:**

* `ooxmlx_<version>_<os>_<arch>.zip`
* `SHA256SUMS.txt`

**Example verification:**

```bash
curl -LO https://github.com/<user>/ooxmlx/releases/download/v1.0.0/ooxmlx_v1.0.0_linux_amd64.zip
curl -LO https://github.com/<user>/ooxmlx/releases/download/v1.0.0/SHA256SUMS.txt
sha256sum -c SHA256SUMS.txt | grep linux_amd64
./ooxmlx --version
```

---

## 🧰 Developer Practices

**Code Standards**

* Go 1.23+
* `CGO_ENABLED=0`
* `go vet` and `go fmt` enforced
* Minimal interface surface (composition over inheritance)
* Prefer list comprehensions (`for range`) and local variable caching for speed

**Testing Practices**

* Must pass race detector.
* 90%+ coverage target for `internal/` packages.
* No file I/O outside `t.TempDir()` in tests.
* Fuzz each parser function with varied seeds.

**Security Practices**

* Mandatory path safety checks for all extracted members.
* No file creation outside destination directory.
* No dynamic imports or unsafe reflection.

---

## 🚀 Future Agents

1. **Parallel Extractor Agent** — Worker pool to parallelize large archive processing.
2. **Checksum Agent** — SHA256 validation of extracted files.
3. **Config Agent** — JSON-based configuration for batch processing.
4. **Golden Test Agent** — Regression tests with stable XML outputs in `testdata`.

---

## ✅ Summary

Each agent in the `ooxmlx` ecosystem has a **single responsibility**, strong **test coverage**, and explicit **build + release contracts**.
The architecture remains CGO-free, dependency-free, and aligned with Go’s best practices for reliability, maintainability, and reproducibility.
