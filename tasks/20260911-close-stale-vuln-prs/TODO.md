# TODO: Close stale vulnerability-fix PRs & prevent recurrence

## Background

`.github/workflows/vulnerability-scan.yml` runs every Monday and, whenever
`govulncheck` reports a problem, creates a brand-new branch
(`fix/vuln-<timestamp>`) and a brand-new PR — without ever checking whether a
PR from a previous run is still open. Since this workflow has never closed
its own prior PRs, 12 duplicate PRs (#27–#38, all titled
"fix: resolve vulnerabilities detected by govulncheck") have piled up since
2026-06-15, all superseded by later runs.

User request (paraphrased after clarification): close the currently-open
stale auto-PRs, and fix the workflow so this doesn't keep happening — the
newest scan's PR should supersede/close the previous one automatically.

## Requirements

- [ ] Close PRs #27–#38 on GitHub (all superseded duplicates), with a comment
      explaining why.
- [ ] Edit `vulnerability-scan.yml` so that, right before it creates a new
      fix PR, it closes any still-open PR from a previous run with the same
      label/title pattern (commented as superseded).
- [ ] Do not touch the "create issue if vulnerabilities found" logic — that
      part already deduplicates correctly (comments on existing issue).
- [ ] Keep the change minimal — no new abstractions, no new workflow files.

## Out of scope

- Not merging or reviewing the *contents* of the old PRs' dependency bumps.
- Not changing the govulncheck / go.mod update logic itself.
