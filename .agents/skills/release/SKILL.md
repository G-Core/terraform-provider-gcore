---
name: release
description: >
  Release workflow for the G-Core/terraform-provider-gcore Terraform provider
  repository. Invoke with /release to discover the open Stainless release PR,
  analyze its diff and changelog, check CI status, merge with user confirmation,
  and generate human-readable release notes on the GitHub Release.
disable-model-invocation: true
allowed-tools: >
  Read,
  Bash(gh pr list --repo G-Core/terraform-provider-gcore *),
  Bash(gh pr view * --repo G-Core/terraform-provider-gcore *),
  Bash(gh pr diff * --repo G-Core/terraform-provider-gcore *),
  Bash(gh pr checks * --repo G-Core/terraform-provider-gcore *),
  Bash(gh pr merge * --repo G-Core/terraform-provider-gcore *),
  Bash(gh release view --repo G-Core/terraform-provider-gcore *),
  Bash(gh release edit * --repo G-Core/terraform-provider-gcore *),
  Bash(sleep *)
---

# Release Skill — G-Core/terraform-provider-gcore

## Constraints

- **Repository**: `G-Core/terraform-provider-gcore` (hardcoded, do not use for
  other repos)
- **Allowed tools**: `gh` CLI (scoped to `G-Core/terraform-provider-gcore`) and
  `Read`
- **Never**: modify source code, force-push, delete branches, or merge without
  explicit user confirmation
- **Release PRs** are created by `stainless-app[bot]` with title
  `release: {version}`

## Workflow

Execute steps 1-6 in order. Present findings at each step before proceeding.

### Step 1 — Discover the Release PR

```bash
gh pr list --repo G-Core/terraform-provider-gcore --state open \
  --app stainless-app --json number,title,author,url,createdAt
```

Find the PR whose title starts with `release: `.

- If **no open release PR** exists, inform the user and stop.
- If found, display: PR number, title (contains version), URL, creation date.

### Step 2 — Analyze the Release

Fetch all three in parallel:

1. PR body containing the auto-generated changelog (Part 2). Parse it to
   extract:
   - Version number and date
   - Breaking changes, Features, Bug Fixes, Refactors, Chores, Documentation

   ```bash
   gh pr view {N} --repo G-Core/terraform-provider-gcore \
     --json body,title,number,url
   ```

2. The actual code diff. Analyze for Terraform provider changes:
   - New/removed/renamed resources and data sources
   - New/changed/removed attributes in resource or data source schemas
   - Changed attribute behaviors (Required, Optional, Computed flags)
   - New plan modifiers or validators
   - Custom code additions (acceptance tests, sweepers, import support)
   - Note: `DeprecationMessage` in schema attributes is informational — the
     attribute still exists. Do **not** surface these as user-facing changes
     in release notes unless the attribute is actually removed in this release
     AND accompanied by a migration path.

   ```bash
   gh pr diff {N} --repo G-Core/terraform-provider-gcore
   ```

3. List of changed files. Use file paths to infer product areas:

   ```bash
   gh pr diff {N} --repo G-Core/terraform-provider-gcore --name-only
   ```

   | File path prefix | Product Area |
   |---|---|
   | `internal/services/cdn_*` | **CDN** |
   | `internal/services/cloud_*` | **Cloud** |
   | `internal/services/dns_*` | **DNS** |
   | `internal/services/fastedge_*` | **FastEdge** |
   | `internal/services/waap_*` | **WAAP** |
   | `internal/services/storage_*` | **Object Storage** |
   | `internal/services/streaming_*` | **Streaming** |
   | `internal/services/iam_*` | **IAM** |
   | `internal/services/security_*` | **DDoS Protection** |
   | Everything else (`internal/apijson/`, `internal/customfield/`, `internal/customvalidator/`, `internal/planmodifier/`, `internal/acctest/`, `internal/sweep/`, `go.mod`, `.github/`, `scripts/`, `docs/`, etc.) | **Other** |

   Each service directory under `internal/services/` maps to a Terraform
   resource or data source. Derive the Terraform resource name by prefixing
   `gcore_` to the directory name (e.g., `cloud_load_balancer/` →
   `gcore_cloud_load_balancer`). Use the full resource name as the sub-area
   header.

   These are heuristics. When a path does not clearly map, use judgment
   based on the directory name and diff context.

   Resources and data sources get **separate sub-area headers**. Data source
   directories typically have plural names (e.g., `cloud_networks/` for the
   list data source vs `cloud_network/` for the resource).

### Step 3 — Check CI Status

```bash
gh pr checks {N} --repo G-Core/terraform-provider-gcore \
  --json name,state,bucket
```

Exit codes: `0` = all pass, `8` = pending, `1` = failure.

| Status | Action |
|---|---|
| exit `0` / all checks pass | Report **CI green**, proceed |
| exit `8` / pending | Warn user checks are running. Ask: wait or proceed? |
| exit `1` / failure | Show failing checks. **Do not offer to merge.** |

### Step 4 — Generate Human-Readable Release Notes

Read `references/release-notes-examples.md` for style reference and examples.

Using the changelog (Step 2.1) and diff analysis (Step 2.2), generate the
release notes following this structure:

#### Disclaimer (always included)

The disclaimer is always prepended at the very top of the release notes, before
Part 1. Insert the actual version number (without the `v` prefix) in both the
prose and the HCL block:

```markdown
> [!warning]
> v2 is a ground-up rewrite of the provider, featuring OpenAPI-spec-driven
> code generation and a move to terraform-plugin-framework under the hood.
>
> This is an **alpha** release, and **breaking changes** are expected.

If you'd like to try it out, pin the provider version **exactly** to
v{VERSION} in your Terraform configuration.

\```hcl
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "{VERSION}"
    }
  }
}
\```

## Release notes
```

#### Part 1 — Human-Readable Summary

After the disclaimer, separated by a `## Release notes` heading, add the
human-readable summary:

```markdown
We're excited to announce version {VERSION}!

### **{Product Area}**

* **`gcore_{resource}`**
  * Added `{attribute}` attribute — {description}
  * ⚠ BREAKING CHANGE: Removed `{attribute}` attribute — {reason}
  * Deprecated `{attribute}` attribute — use `{alternative}` instead
  * Fixed {description} — {detail}
```

#### Part 1 Rules

- **Group by product area** (alphabetically: CDN, Cloud, DDoS Protection,
  DNS, FastEdge, IAM, Object Storage, Streaming, WAAP). Place **Other** last.
  Only include areas that have changes.
- **Within each area, group by Terraform resource name** alphabetically
  (e.g., `gcore_cloud_floating_ip`, `gcore_cloud_instance`,
  `gcore_cloud_load_balancer`). Each resource and data source gets its own
  sub-area header.
- **Sub-area headers** use bold + backtick: `* **\`gcore_cloud_network\`**`
- **Attribute names** use snake_case Terraform names in backticks (e.g.,
  `create_router`, `admin_state_up`) — what users see in their `.tf` files.
- **Breaking changes** get `⚠ BREAKING CHANGE:` prefix inline. Always specify
  what was removed or changed and why.
- **Deprecations**: ``Deprecated `{attribute}` attribute — use `{alternative}`
  instead``
- **Additions**: ``Added `{attribute}` attribute — {description}``
- **Fixes**: `Fixed {what} — {detail}`
- **No commit hashes or links** in Part 1 (those are in Part 2).
- **Do not copy** the auto-generated changelog verbatim. Aggregate related
  changes. Skip noise.
- **Omit** `codegen metadata`, `aggregated API specs update`,
  `codegen related update`, `import path fix`, and `stale module path` entries
  unless they introduce a specific user-visible change visible in the diff.

Display the generated disclaimer + Part 1 to the user. Ask if they want to
edit or approve.

### Step 5 — Merge the Release PR

Present to the user:
1. CI status (from Step 3)
2. Generated release notes preview (disclaimer + Part 1 from Step 4)
3. Auto-generated changelog (Part 2 from PR body)
4. Recommended merge method: **rebase**

**Ask for explicit confirmation before merging.**

If user declines or wants changes, return to Step 4.

Once confirmed, merge via Bash:
```bash
gh pr merge {PR_NUMBER} --repo G-Core/terraform-provider-gcore --rebase
```

If merge fails, report the error and stop.

### Step 6 — Update the GitHub Release

After merge, `stainless-app[bot]` auto-creates a GitHub Release. GoReleaser
then builds and attaches provider binaries for all platforms (this is handled
by CI, not this skill).

1. Fetch the latest release. Verify `tagName` matches expected version
   `v{VERSION}`. If not found, `sleep 10` and retry once.

   ```bash
   gh release view --repo G-Core/terraform-provider-gcore \
     --json tagName,body,url
   ```

2. Build the final release body by combining the disclaimer, Part 1, and the
   existing Part 2:

   ```
   {Disclaimer}

   ## Release notes

   {Part 1 — human-readable summary}


   {Part 2 — auto-generated changelog already in release body}
   ```

3. Update the release via Bash:
   ```bash
   gh release edit v{VERSION} \
     --repo G-Core/terraform-provider-gcore \
     --notes "$(cat <<'RELEASE_EOF'
   {combined release notes}
   RELEASE_EOF
   )"
   ```

4. Display the release URL and confirm completion.

## Failure Modes

| Situation | Action |
|---|---|
| No open release PR | Inform user, stop |
| CI failing | Show failures, do not merge |
| CI pending | Warn, ask user preference |
| Merge conflict | Report, suggest manual resolution |
| Merge fails | Report error, stop |
| Release not found after merge | Retry once after 10s, then report |
| `gh` CLI not authenticated | Report, suggest `gh auth login` |
