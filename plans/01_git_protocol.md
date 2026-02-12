# Agent 1: Core Git Engine & Protocol (Rust)

**Mission:** Implement the core Git connectivity protocols to allow users to `git push` and `git clone` repositories reliably. This is the foundation of the platform.

**Status:** Completed (Smart HTTP & Hooks)

## Goals

1.  **Implement Smart HTTP Protocol** (Completed)
    *   Endpoint: `POST /repo/{name}/git-upload-pack` (Fetch/Clone) - Done
    *   Endpoint: `POST /repo/{name}/git-receive-pack` (Push) - Done
    *   Ensure correct Content-Type headers (`application/x-git-upload-pack-result`, etc.) are handled. - Done
    *   Implement authentication verification (Basic Auth or Bearer Token) before processing requests. - Done (JWT via Basic Auth)

2.  **Implement SSH Transport** (Deferred)
    *   Develop an SSH server (using `russh` or similar crate) that authenticates users via public keys stored in the DB.
    *   Forward Git commands to the underlying `git` binary or `libgit2`.

3.  **Git Hooks & Verification** (Completed)
    *   Implement `pre-receive` hooks to enforce branch protection rules (e.g., no force push to main). - Done
    *   Implement `update` hooks to validate commit signatures (MIP integration point). - Planned for MIP agent.
    *   Implement `post-receive` hooks to trigger webhooks and CI/CD events. - Planned for CI/CD agent.

4.  **Storage Optimization** (Deferred)
    *   Implement `git gc` (Garbage Collection) scheduling.
    *   Implement object deduplication strategies for forks (alternates).

## Tech Stack
*   **Language:** Rust (Actix-Web)
*   **Libraries:** `git2` (libgit2), `tokio` (Async runtime), `jsonwebtoken`

## Verification
*   User can run `git clone http://localhost:8081/repo/my-repo.git` successfully. (Verified)
*   User can run `git push origin main` successfully. (Verified)
*   Pushing to a protected branch is rejected with a meaningful error message. (Verified)
