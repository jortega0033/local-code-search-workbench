# Implementer map

Status: **planned, not implemented.** This page maps the boundaries that [PROJECT-PLAN.md](PROJECT-PLAN.md) requires onto a module layout so that parallel tickets do not collide. It adds no requirement; the plan wins where they differ. The names are the orchestrator's proposal, fixed by the scaffold; a ticket that changes one corrects this page in the same pull request.

## Planned module map

| Boundary | Module | Owns | Parent ticket |
|---|---|---|---|
| Entry point | `cmd/codesearch` | Subcommands `serve`, `generate`, `verify` | #2 (`generate`), #5 (`serve`), #7 (`verify`) |
| Corpus and snapshot | `internal/corpus` | Walk, admission, skip counts, language map, line starts, immutable snapshot, generator | #2 |
| Query and reference scan | `internal/query` | Modes, validation, filters, line matching, segments, caps, force scan | #3 |
| Trigram index | `internal/index` | Postings, prefix rule, intersection, budget, strategy counters | #4 |
| HTTP | `internal/httpapi` | Routes, guards, versions, rebuild, queue, cancellation, static assets | #5 |
| Browser assets | `web/ (embedded)` | Query, filters, results, context, status | #6 |
| Verification | `internal/verify` | 60-query matrix, parity, grep comparison, timing report | #7 |

Dependency direction: `httpapi` depends on `query`, `index` and `corpus`; `index` and `query` on `corpus`, never on each other; `verify` on all of them. `corpus` imports nothing from the module.

## Freeze points

Parallel tickets rely on these shapes. A ticket that changes a frozen shape updates this page and every dependent ticket in the same pull request.

- **Snapshot contract** (file ID, relative path, language, bytes, line starts, version): ticket #2; #3 and #4 build on it in parallel.
- **Result row and segment shape**, and the candidate set shape: agreed between #3 and #4 before #5 starts.
- **HTTP JSON** (collection, reindex, search, context, error causes): #5; read by #6 and #7.

## What is verified where

Go is verified in CI only (see [CLAUDE.md](../CLAUDE.md)), on Linux, where symlink, FIFO and GNU grep cases can run. Manual browser steps are owner checks in Chrome.
