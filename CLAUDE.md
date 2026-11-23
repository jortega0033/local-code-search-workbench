# CLAUDE.md

Working instructions for this repository. Read this file first, then the documents it links to. It does not repeat them.

## What this is

A read-only search tool for a developer with a local copy of a codebase: build a bounded snapshot, search it with a literal or a regular expression, filter by path and language, and read exact highlights in context. Go standard library only, `127.0.0.1` only, browser assets embedded, every indexed result checkable against a full scan.

**Nothing is implemented yet.** The repository holds the plan, the delivery tooling and the backlog in epic #1. Commands marked planned below do not work yet. Do not describe a planned thing as working, and do not claim a passing baseline: there is none.

## Sources of truth

| Topic | Document |
|---|---|
| Product scope, design, fixtures with expected values | [docs/PROJECT-PLAN.md](docs/PROJECT-PLAN.md) |
| Planned module map, freeze points, planned versus implemented | [docs/architecture.md](docs/architecture.md) |
| Backlog | Epic issue #1 and its child tickets on GitHub. The epic's table is authoritative for order, dependencies and estimates; tickets state their own owner-decision defaults. |
| Loop setup and its limits | [LOOP.md](LOOP.md), [loop-constraints.md](loop-constraints.md) |
| Project subagents, origin and local changes | [.claude/agency-agents/README.md](.claude/agency-agents/README.md) |
| Ticket workflow and the step at which each role runs | [.claude/skills/implement-ticket/SKILL.md](.claude/skills/implement-ticket/SKILL.md), section "Agent roles" below |
| How a ticket goes from the queue to a ready pull request, who may do what | "Delivery workflow" below |
| UI and copy rules, design tokens | [docs/ui-quality.md](docs/ui-quality.md), [docs/design-tokens.md](docs/design-tokens.md) |
| Decisions the owner made, with dates | [docs/owner-decisions.md](docs/owner-decisions.md) |

The plan wins for product design; the epic and tickets win for order and scope. A ticket that changes a contract fixes the plan in the same pull request.

## Stack and constraints

- Go module, **no external dependencies**: standard library only, in the application and in the tests.
- One executable; HTML, CSS and JavaScript are embedded with `go:embed`. The server binds loopback.
- Read-only: one corpus root from the command line, nothing written inside it, no path taken from the browser.
- Offline: same-origin assets only. Fixtures are original and generated from a seed.
- Limits: 2,000 files, 64 MiB accepted, 4 MiB per file, 256 KiB per line, 512-byte patterns, 200 rows, 20 ranges per line.
- Editing, remote crawling, cross-line matches, live updates, several users and the 1 GiB milestone are out of scope.
- A runtime version change, a dependency, an external request or a non-loopback listener needs the owner's approval first, not a pull request.

## Commands

Go is not installed on the owner's PC and nothing is installed globally. **Go checks run in CI only.** Run `go version` first: if it is missing, do not install it, and do not claim you ran any command below.

| Purpose | Command | Status |
|---|---|---|
| Build | `go build -o codesearch ./cmd/codesearch` | planned: real with the scaffold |
| Format check (empty output is good) | `gofmt -l .` | planned: real with the scaffold |
| Static check | `go vet ./...` | planned: real with the scaffold |
| Tests | `go test ./...` | planned: real with the scaffold |
| Generate the corpus | `./codesearch generate --seed 1 --out ./corpus` | planned: #2 |
| Run the app | `./codesearch serve --root ./corpus` | planned: #5 |
| Parity and timing report | `./codesearch verify` | planned: #7 |

## Architecture boundaries

The planned module map is in [docs/architecture.md](docs/architecture.md). Invariants every change must keep (the plan is the source):

- Corpus admission, query evaluation, trigram index, HTTP and page stay separate. The index yields candidates only and never decides a match.
- Files are opened through the root-scoped file API (`os.Root`), after a lexical walk that rejects a symlink root and skips symlinks, dot directories, VCS metadata, non-regular files and the documented generated folders.
- Skip files over 4 MiB, with NUL, invalid UTF-8 or a line over 256 KiB, counted per cause. Hard maxima: 2,000 files, 64 MiB accepted bytes, a posting budget. An overflow fails clearly and publishes nothing.
- A published snapshot is immutable and pins the captured bytes. One rebuild at a time; a failed rebuild keeps the previous snapshot; an old version gets a typed stale answer.
- Patterns: 1 to 512 UTF-8 bytes, literal or regex mode, compiled once per request. Invalid patterns and filters never start a scan.
- A line ends at LF, a preceding CR is not matched text, a final line without LF counts, a final LF adds no line. Spans are file-absolute, zero-based, half-open bytes, in the engine's non-overlapping order.
- Byte trigrams come only from `LiteralPrefix()` of at least three bytes; otherwise the strategy is a scan. Path filters match whole segments: `src` excludes `src-old`.
- Order is normalised path, then line. At most 200 rows with a probe for the 201st, 20 ranges per line with an omitted flag, and an explicit cause for every incomplete answer.
- The server binds loopback, checks host and origin, accepts writes as same-origin JSON only, takes no filesystem path from the browser and never returns the absolute root.
- One active search and one queued request at most; cancellation is checked between bounded lines or files; a time limit gives an incomplete answer with its own cause.
- Filenames and snippets are text nodes. Highlights come from the server's `match` and `plain` segments. Same-origin assets only.
- Query, filters, results and context work by keyboard. Strategy, skips, truncation, build state and errors are text, a zero-width match has a visible marker, and the last good results stay with their version.
- The reference evaluator makes no index call. GNU grep is compared on path and line only, on generated LF-only ASCII files and the documented subset; a missing grep is reported, not passed.
- Performance numbers are measured on a named machine, toolchain and corpus hash, with the repetition count. No number is written before it was measured.

## Git workflow

- `main` is the default branch. Work on one branch per ticket: `<issue-number>-<short-slug>`. Never commit to `main`.
- Small commits with imperative subjects and a type prefix (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).
- Order for every ticket: local evidence, Code Reviewer, Reality Checker, then a draft pull request that says `Part of epic #1` and `Closes #<ticket>`, then CI on that draft, then `gh pr ready` only when CI is green.
- Do not push or open a pull request unless the owner asked in the current session; starting `/next-ticket` is that request for one ticket. Only the owner merges and releases. Do not push to `main`, force-push, or merge your own pull request.
- Never commit generated corpora (`corpus/`), `.env*`, or build output.
- Do not change `.claude/agents/`, `LOOP.md`, `loop-constraints.md` or `loop-budget.md` as part of a feature ticket.
- Open owner decisions on a ticket are settled with the ticket's stated default. Where a ticket states none, take the conservative option and say so in the pull request.

## Evidence before a ticket is called done

Paste real output into the pull request; never summarise a run you did not make.

1. For each acceptance criterion of the issue: the name of the test that proves it, or the exact step and what was observed.
2. Go results come from CI, and green CI on the current diff is a condition for `gh pr ready`: the link to the CI run on the draft pull request with the Go version and the results of `gofmt -l .`, `go vet ./...` and `go test ./...`. If Go is not available locally, say so.
3. For a search or index change: the fixture rows with path, line and byte spans, and the indexed against scan comparison.
4. For an HTTP change: one request and response pair per success and per error cause (CI test output counts).
5. For a UI change: the browser and version, the steps and what was observed, and the manual steps left as owner checks.
6. `git status` shows no stray files, and no dependency, external URL or non-loopback listener was added.
7. Anything not verified, said plainly.

## Agent roles

Project subagents are in `.claude/agents/`: eight Agency Agents roles and `loop-verifier`, which belongs to the Loop setup. Each Agency Agents role file starts with a section "In this repository", written for this project. That section and this file override the generic text below it. Sources, licence, local changes and the name collision with a user-level install are in [.claude/agency-agents/README.md](.claude/agency-agents/README.md).

When asked to implement, continue or finish a GitHub issue, follow the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Claude Code loads it on its own for ticket work; `/implement-ticket <issue number>` starts it by hand. Delegate to the roles without waiting to be asked: each role's `description` ends with a "Use proactively ..." sentence, which is what makes Claude Code delegate. The playbook gives the exact trigger of each role.

| When | Roles | Access |
|---|---|---|
| Before coding | Software Architect | `Read, Grep, Glob` |
| While coding | Backend Architect, Frontend Developer | can edit |
| After coding, when the trigger matches | Application Security Engineer, Accessibility Auditor | `Read, Grep, Glob` |
| After coding, when the trigger matches | API Tester | can edit |
| Gate on every ticket | Code Reviewer on the full diff, until no material finding is left | `Read, Grep, Glob` |
| Acceptance gate on every ticket, after Code Reviewer | Reality Checker on each acceptance criterion and the recorded evidence, and that CI is configured to run the checks (no CI result is needed yet) | `Read, Grep, Glob` |

- The two gates run on every ticket, whatever its size. A ticket is not called done before both have passed.
- Roles limited to `Read, Grep, Glob` cannot run commands. Give them the diff, the acceptance criteria and whatever command output exists.
- The session or role that wrote a change does not count as its reviewer.
- If a role is missing or cannot be delegated to, the session does the before-coding, while-coding and specialist steps itself and says so. It cannot pass the two gates on its own change: do not push or open the pull request. Stop, report that the gates could not run, and leave the run for the owner. This overrides the playbook's sentence "do that step yourself" for the two gates.
- In the pull request description, list under "Agent roles run" each role that ran, what it was asked and its verdict, and any step done without a role.

## Delivery workflow

How a ticket reaches the owner. Each piece is described in its own file; this part only links them.

- **Start:** the owner runs `/next-ticket` ([`.claude/skills/next-ticket/SKILL.md`](.claude/skills/next-ticket/SKILL.md)); it picks one leaf ticket (never a `tracking` parent or a ticket with open sub-issues; the last sub-issue's pull request also closes the parent), creates the branch and worktree, and runs the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Nothing is scheduled.
- **Gate hook:** [`.claude/hooks/role_gate.py`](.claude/hooks/role_gate.py) blocks `gh pr create` and `gh pr ready` until Code Reviewer and Reality Checker ran after the last code edit. `ROLE_GATE=off` in the environment Claude Code is started from switches it off for that session; only the owner does that.
- **Permissions:** [`.claude/settings.json`](.claude/settings.json). Agent configuration edits ask in the main checkout and in ticket worktrees; CI, `CLAUDE.md`, `AGENTS.md` and Loop-file edits ask only in the main checkout; in a ticket worktree they are limited by the playbook's scope rule and reviewed by the owner at merge (shell-made edits are not covered, so a ticket never changes agent configuration); force pushes, pushes to `main`, `gh pr merge`, releases and reads of `.env` are denied.
- **CI:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml). Every job must be green on the current diff before a pull request is marked ready. The agent does not poll: it ends its turn with "Draft pull request #N is waiting for CI" and reads `gh pr checks` once when woken. Reviewers read the diff and evidence from `.claude/state/review/`. An acceptance step that needs a person is listed under "Owner checks" in the pull request.
- **Pull request:** [`.github/pull_request_template.md`](.github/pull_request_template.md). Tickets use [`.github/ISSUE_TEMPLATE/implementation-ticket.md`](.github/ISSUE_TEMPLATE/implementation-ticket.md).
- **Labels:** `status:ready`, `status:blocked`, `status:in-progress`, `needs-decision`, `tracking` (parent container, never worked; its sub-issues are the work items).
- **Changelog:** every ticket pull request adds one line under Unreleased in [CHANGELOG.md](CHANGELOG.md).
- **UI and copy:** [docs/ui-quality.md](docs/ui-quality.md); the `ui-ux-pro-max` and `unslop` skills are in `.claude/skills/`. Colours come from the shared design tokens in `web/tokens.css` ([docs/design-tokens.md](docs/design-tokens.md)); new UI adds no colour literals.
- **Decisions:** [docs/owner-decisions.md](docs/owner-decisions.md). Security reports: [SECURITY.md](SECURITY.md).
- **Merge and release:** only the owner merges and releases ([docs/RELEASING.md](docs/RELEASING.md)). An agent stops at a ready pull request.

## Loop

Issue triage runs at L1: report-only, started by hand, no scheduler, no GitHub writes. See [LOOP.md](LOOP.md).
