# Local Code Search Workbench

> This document specifies phase 1 of the program in [#39](https://github.com/jortega0033/local-code-search-workbench/issues/39). Hour figures in it are early estimates for that phase; its open tickets add up to 72 to 100 hours (2026-11-24). Work that it defers is planned in the later phases only where the program names it.

## Product requirements

A developer with a local copy of a codebase can find matching lines across the collection, narrow them by language or path and inspect the exact matches in context. Chrome shows the query, indexed and skipped file counts, search strategy and whether results are complete. Opening a result shows neighboring lines from the same indexed snapshot. The project starts with an original synthetic corpus.

For the first useful slice, point the server at a fixture root, build a collection and run a literal search from Chrome. Pick a result to see the matching excerpt, filename and one-based line number. The first release adds regular expressions, path and language filters, UTF-8 highlighting, reindexing and an independent comparison against a force-scan path. It runs with the installed Go compiler and Chrome, served on loopback with no network downloads, service accounts or editor framework.

### Success criteria

- Search a documented 2,000-file corpus within a 64 MiB byte budget. Report file counts, skip reasons, index build time and memory observations. Complete the scripted offline browser walkthrough on a fresh checkout.
- Run 60 checked-in query descriptions against a deterministic generated corpus. For complete searches, compare ordered file paths, lines, and nonempty byte spans between the index and the full scan. Independently assert the hand-calculated fixture below. Record whether the planner uses grams or scans and how many files it inspects.
- Exercise invalid patterns and filters, short patterns and empty highlights, UTF-8 and CRLF, refresh failures, and result caps. Every cap or failure has a visible explanation and leaves the previous successful view usable.
- Produce a repeatable report of cold build time, warm query p50 and p95, indexed-to-scan candidate counts, and elapsed client search time for selective and broad queries. Record measured values alongside hardware, toolchain and corpus hash. Performance numbers are measurements to collect during implementation.

This release is read-only. Source editing, remote crawling, cross-line regex, instant live updates and multi-user access belong outside the first release. A later milestone grows to a 1 GiB corpus, 200 public-subset regex queries and a 100 ms warm query target per query. That milestone needs a larger-storage indexing design, stable benchmark conditions and an explicit output contract. The timing target applies to returning the first 200 rows; check complete match-line parity in a separate untimed run. Include broad-query and scan-fallback timing and report misses explicitly.

## Technical design

### Local runtime and snapshot

Use Go 1.24 or later, standard library only. Check the installed compiler and bundled browser first. Serve HTML, CSS and vanilla JavaScript from the binary using `go:embed`, with no runtime assets outside loopback. The CLI accepts one corpus root and a loopback bind address. The browser cannot supply filesystem paths or read raw files by URL. Host checks accept the configured loopback origin, mutating endpoints require a same-origin JSON request, and unrelated origins receive no CORS grant. Use plain-text DOM nodes for filenames and snippets.

Walk the root lexically, reject an initial symlink root, and skip symbolic links, dot directories, VCS metadata, unsupported file kinds and documented generated folders. Open accepted relative names through the root-scoped file API. Before reading, verify the opened file is regular. Skip files over 4 MiB, files containing NUL or invalid UTF-8 and files with any line over 256 KiB. Set hard maxima of 2,000 files, 64 MiB accepted bytes and a separate bounded index-memory budget based on stored posting entries. For a collection or posting-budget overflow, return a clear error without publishing a partial location index. Use approved original or openly licensed input and fixture files. [Go os.Root](https://pkg.go.dev/os#Root)

Store preflight counts, relative slash-normalized paths, mapped language, byte contents and line starts in an immutable in-memory collection. A published snapshot pins the captured bytes that power results and excerpts, even if the disk copy changes. Build one replacement at a time. A failed rebuild leaves the previous snapshot available. During rebuild, both snapshots may temporarily occupy memory, so bound peak allocation before allocating replacement data. A click on an old result after the snapshot is replaced returns a stale-version message and a refresh action. Restart discards the index and rebuilds it from the supplied root.

### Queries and candidate selection

Provide explicit literal and regex modes. Literal mode escapes regex syntax. Regex mode uses Go’s standard-library engine and its documented syntax. Compile outside rendering and reuse the compiled expression for all lines in the request. Cap the submitted pattern at 512 UTF-8 bytes, reject empty queries, and return syntax errors without running a scan. Go’s regular-expression package documents linear-time operation in input length. Caps on corpus size, requests and result construction also bound work from valid broad queries. [Go regexp package](https://pkg.go.dev/regexp)

A line ends at LF, with a preceding CR excluded from that line’s match text. A final line without LF is searchable. Empty files have zero searchable lines, consecutive LF bytes delimit blank lines and a final LF creates no additional record. Match one line at a time using the engine’s non-overlapping occurrence order. Report file-absolute, zero-based, half-open byte spans. Return original text as safe text segments labeled `match` or `plain`, plus byte offsets, for correct browser highlights. For a zero-width occurrence show its line and visible marker beside the excerpt, with no highlighted character. Admit one active search and one queued request at most. Cancel a superseded browser request and check server cancellation between bounded lines or files.

Index each distinct overlapping UTF-8 **byte** trigram in a file into a sorted list of file IDs. For a normal search call `LiteralPrefix()` on the compiled regex. A documented required prefix of at least three bytes makes every overlapping prefix trigram necessary. Intersect their sorted posting lists, then apply the original regex to eligible candidate lines. A missing necessary posting proves emptiness. When the engine provides no usable prefix, use a full scan within the filters. In either case, intersect candidate files with exact extension-language filters and normalized segment-prefix paths such as `src/ui`. A path prefix `src` includes `src/a.go` and excludes `src-old/b.go`. Fail closed on an invalid filter. Force-scan bypasses postings for the same input and snapshot. [Go Regexp.LiteralPrefix](https://pkg.go.dev/regexp#Regexp.LiteralPrefix)

Results are sorted by normalized path and line number. Return a line reference, up to 20 ranges per line, and flags for omitted ranges. Return up to 200 result rows and probe for a 201st row to distinguish exactly full from truncated. A server scan time limit returns an explicitly incomplete response with its own cause. Skip counts refer to corpus admission, not search truncation. Limit excerpt length on UTF-8 boundaries and use a context view for long lines. Never expose the absolute corpus root in the browser, error objects or JSON responses.

### Interfaces and checks

- `GET /api/collection` returns status, version, indexed size and skip groups. `POST /api/reindex` schedules a single build and returns the status of that attempt. The UI reads the existing version until publication succeeds.
- `POST /api/search` receives the version, mode, pattern and filters. It returns rows, timings, a strategy value of `trigram`, `scan` or `empty`, files examined and truncation flags. Include the typed invalid-pattern and stale-version cases in the contract.
- `GET /api/context` accepts the returned version, internal file ID and line, using a bounded context window. The server only reads from the active snapshot. The UI treats content as text and includes keyboard navigation and a readable empty state.
- Add a direct evaluator for the query-admission and line-byte contracts using manually specified results. In addition, use a scan-only oracle that evaluates each admitted line without index calls. For an independent executable reference, compare **matching path/line pairs** against GNU grep `-n -E` on generated LF-only ASCII files and an explicitly documented common regex subset. Span parity follows the application contract rather than the grep comparison. Treat absence of an installed grep as a separately reported verifier prerequisite. [GNU grep manual](https://www.gnu.org/software/grep/manual/grep.html)

### Small verification fixture

Use the following byte strings, with the escape sequences denoting the actual CR/LF bytes and `é` encoded in UTF-8.

| Path | Contents |
| --- | --- |
| `src/a.go` | `alpha\ncafé todo\nalpha alpha\n` |
| `src/b.py` | `beta\r\nTODo\n` |
| `src-old/guide.txt` | `todo\n` |

Literal `alpha` returns `src/a.go` line 1 span `[0,5)` and line 3 spans `[17,22)` and `[23,28)`. Literal `é` returns `src/a.go` line 2 span `[9,11)` and uses scan because the prefix is two bytes. Regex `[Tt][Oo][Dd][Oo]` with prefix filter `src` returns `src/a.go` line 2 span `[12,16)` and `src/b.py` line 2 span `[6,10)`, in that order. The same query without the prefix includes the guide first because `src-old/` sorts before `src/` in byte order. Regex `^beta$` finds `src/b.py` line 1 span `[0,4)`. Add separate zero-width, reindex and rejected symlink fixtures.

## Epic · Search and verify a bounded local source collection

Deliver the first offline release and its reproducible comparison report. The complete scope includes discovery, index and scan correctness, a local API, a Chrome interface and measured verification. The six tickets total **17–25 estimated engineering hours** including implementation and checks. Start CS-01 alone. CS-02 and CS-03 can run in parallel against the same immutable snapshot and test contracts. CS-04 integrates them. CS-05 connects the interface, and CS-06 verifies the completed offline workflow.

### CS-01 · Build the original corpus generator and snapshot contract

**Estimated effort** 2–3 hours. **Dependencies** None.

Define relative file IDs, language mapping, immutable content and stable line references. Add the tiny fixture and a seeded multi-file generator. Implement corpus reads and per-cause counts in the future build.

**Acceptance checks**

- Independent tests verify lexically ordered paths, the exact tiny fixture bytes, LF and CRLF line starts, and the same corpus hash from two generator runs with identical inputs.
- Cover a dot directory, ordinary symlink, escaping symlink, non-regular entry, invalid UTF-8, NUL, oversized single file and oversized line. Skip these according to the documented rules. Collection overflow fails without publishing the candidate.
- Demonstrate that the generated corpus is original and reconstructible offline from its seed.

### CS-02 · Define query behavior with a reference scan

**Estimated effort** 3–4 hours. **Dependencies** CS-01.

Build the reference evaluator and its display segments, typed validation and exact filtering. All line matching uses captured bytes.

**Acceptance checks**

- The tiny fixture produces every stated path, line and byte span in order. Include a blank line, final unterminated line, capture-free zero-width result and a line with more than 20 ranges.
- Show that path filter `src` excludes `src-old`. Verify language and path combination, escaped literal syntax, malformed patterns, just-under and over-limit pattern lengths, and the difference between 200 and at least 201 matching rows.
- The returned UTF-8 fragments reconstruct the original line where no excerpt cap applies. Invalid filters or queries never start evaluation.

### CS-03 · Build trigram postings and a conservative query planner

**Estimated effort** 3–5 hours. **Dependencies** CS-01.

Build the in-memory posting store, prefix eligibility rule, ordered intersection and strategy counters. Produce candidates independently of text matching.

**Acceptance checks**

- Verify exact candidate IDs for a prefix with repeated grams, a missing posting, a UTF-8 prefix spanning multiple bytes and a prefix present in a different order from the pattern. The last can be an extra candidate and never a returned false match.
- Property-check, on independently generated fixture ranges, that each prefix actually occurring anywhere in a file leaves that file in the candidate set.
- Queries with insufficient prefixes, case-alternation patterns and a short prefix enter scan strategy. A failed allocation budget leaves no partially usable index.

### CS-04 · Serve queries, context and versioned rebuilds

**Estimated effort** 3–4 hours. **Dependencies** CS-02 and CS-03.

Connect the local server to both evaluation paths. Pin the snapshot during a search and make new requests choose the current version. Return explicit status for queue saturation, time-budget truncation and stale context.

**Acceptance checks**

- Contract tests compare candidate-filtered and force-scan responses on the tiny corpus, including a path filter, invalid regex and an empty result.
- Replace a disk file, confirm an existing version’s excerpt stays unchanged, rebuild, then confirm the new version returns new bytes and the old version yields a stale error. An invalid replacement keeps the old version searchable.
- Test host and origin checks, malformed IDs, simultaneous builds and cancellation. Verify no endpoint accepts a client-supplied file path or returns the actual absolute root.

### CS-05 · Build the Chrome search and excerpt workflow

**Estimated effort** 3–4 hours. **Dependencies** CS-04.

Make query entry, filter selection, highlights and context inspection usable with just a keyboard. Display strategy, skipped files, incomplete results, build status and errors. Keep last good results visible with clear version status.

**Acceptance checks**

- With external requests blocked, perform literal search, regex search and filter toggling from the keyboard. Open the `café todo` fixture and inspect the exact visible `é` and `todo` highlights.
- Verify that fixture strings resembling HTML appear as text with no executable elements. Validate short, long and zero-width matches visually.
- Exercise reindex success, reindex failure, stale context, invalid regex, empty results, error recovery and a result-cap notice.

### CS-06 · Verify parity and report measured search performance

**Estimated effort** 3–5 hours. **Dependencies** CS-05.

Build the release checks, 60-query matrix and reproducible benchmark procedure. Use a fresh offline run and publish evidence with the eventual implementation.

**Acceptance checks**

- Compare full ordered indexed and force-scan rows and spans on the generated corpus without presentation truncation. Cover literals, regex variants, filter combinations, no-prefix scans and Unicode. Compare common-subset ASCII matching lines against GNU grep independently.
- Record corpus SHA-256, toolchain, environment, cold build, peak observed memory, and per-query p50 and p95 under a documented warm-run repetition count. List candidates and scanned files for selective and broad queries. Treat parity failures and incomplete matrix execution as release blockers.
- From a clean checkout with registry access disabled, run the checked-in Go tests, launch the app and capture an actual browser walkthrough against the bundled fixtures. Record commands, outputs, failure cases and next milestone constraints.
