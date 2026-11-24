# Local Code Search Workbench

An offline source search workbench with a trigram index, verified regex matches and highlighted browser results.

Only the scaffold exists (a Go module with empty packages, smoke tests and CI); the product behaviour is planned. The scope, architecture and acceptance criteria are in [the project plan](docs/PROJECT-PLAN.md). The plan includes a reproducible verification strategy for the future build.

- [Product requirements](docs/PROJECT-PLAN.md#product-requirements)
- [Technical design](docs/PROJECT-PLAN.md#technical-design)
- [Implementation backlog and epic](https://github.com/jortega0033/local-code-search-workbench/issues)

Estimated implementation and verification effort is 17–25 engineering hours.

## Program goal

Build a local code intelligence engine: exact text, regular-expression, structural and semantic search over billions of lines and full history, with precise definitions and references, from one offline Go binary and no search, parser or regex library.

The program is planned ([#39](https://github.com/jortega0033/local-code-search-workbench/issues/39)). Phase 1 is the current epic [#1](https://github.com/jortega0033/local-code-search-workbench/issues/1) and the only phase with tickets.

**Score.** Share of queries in a fixed, generated set where results equal a full scan or an independent compiler's answer, with the median query time on a one-billion-line corpus.

| Phase | Content | Hours (rough) | Issue |
|---|---|---|---|
| 1 | The current epic and its tickets | 72 to 100 | [#1](https://github.com/jortega0033/local-code-search-workbench/issues/1) |
| 2 | Scale | 130 to 220 | [#40](https://github.com/jortega0033/local-code-search-workbench/issues/40) |
| 3 | Query engine | 140 to 240 | [#41](https://github.com/jortega0033/local-code-search-workbench/issues/41) |
| 4 | Structural search | 160 to 270 | [#42](https://github.com/jortega0033/local-code-search-workbench/issues/42) |
| 5 | Code intelligence | 180 to 310 | [#43](https://github.com/jortega0033/local-code-search-workbench/issues/43) |
| 6 | Semantic search and workbench | 120 to 200 | [#44](https://github.com/jortega0033/local-code-search-workbench/issues/44) |

About 800 to 1340 engineer hours in total: rough estimates, not measurements. The numbers in the score and the finish lines are targets until a run shows them. Only phase 1 has tickets to work; the sub-issues of a later phase are written when it starts. No single session is expected to reach the end.

## License

[MIT](LICENSE) © 2026 Jake Ortega.
