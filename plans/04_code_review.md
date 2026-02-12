# Agent 4: Code Review & Collaboration (Go/Rust)

**Mission:** Build the collaboration engine for teams, specifically centered around Pull Requests, Code Review, and sophisticated merge strategies.

## Goals

1.  **Pull Requests (PRs)**
    *   Endpoint: `POST /repos/{owner}/{repo}/pulls`
    *   Track diffs between `head` and `base` branches.
    *   Implement "Draft" state for work-in-progress.
    *   Implement PR templates (e.g., `PULL_REQUEST_TEMPLATE.md`).

2.  **Diff Generation & Syntax Highlighting**
    *   Utilize `smsly-git` (Rust) to generate side-by-side or unified diffs for large files efficiently.
    *   Cache diff results (Redis) to avoid recalculation on every page load.
    *   Handle merge conflicts detection (`git merge-tree`).

3.  **Code Review**
    *   Implement inline commenting on diffs (`POST /repos/{owner}/{repo}/pulls/{number}/comments`).
    *   Support threaded conversations (Resolve/Unresolve).
    *   Implement "Request Changes" vs "Approve" reviews.
    *   (Optional) Suggest code changes directly in the comment.

4.  **Merge Strategies**
    *   Implement `Merge Commit` (Standard).
    *   Implement `Squash and Merge` (Combine all commits).
    *   Implement `Rebase and Merge` (Linear history).
    *   Enforce branch protection rules (e.g., "Require passing checks", "Require approval").

5.  **Conflict Resolution**
    *   Detect conflicts automatically.
    *   Provide a web UI to resolve simple conflicts (optional, advanced).

## Tech Stack
*   **Language:** Go (API Logic), Rust (Git Operations/Diffing)
*   **Database:** PostgreSQL (PR metadata, comments)
*   **Libraries:** `git2` (for diffing/merging)

## Verification
*   User can open a PR from a feature branch to `main`.
*   User can comment on a specific line of code in the diff.
*   Merging blocked if branch protection rules are not met.
*   Squash merge results in a single commit on target branch.
