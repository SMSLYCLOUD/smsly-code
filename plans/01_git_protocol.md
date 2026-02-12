# Agent 1: Core Git Engine & Protocol (Rust)

**Mission:** Implement the core Git connectivity protocols to allow users to `git push` and `git clone` repositories reliably. This is the foundation of the platform.

## Goals

1.  **Implement Smart HTTP Protocol**
    *   Endpoint: `POST /repo/{name}/git-upload-pack` (Fetch/Clone)
    *   Endpoint: `POST /repo/{name}/git-receive-pack` (Push)
    *   Ensure correct Content-Type headers (`application/x-git-upload-pack-result`, etc.) are handled.
    *   Implement authentication verification (Basic Auth or Bearer Token) before processing requests.

2.  **Implement SSH Transport**
    *   Develop an SSH server (using `russh` or similar crate) that authenticates users via public keys stored in the DB.
    *   Forward Git commands to the underlying `git` binary or `libgit2`.

3.  **Git Hooks & Verification**
    *   Implement `pre-receive` hooks to enforce branch protection rules (e.g., no force push to main).
    *   Implement `update` hooks to validate commit signatures (MIP integration point).
    *   Implement `post-receive` hooks to trigger webhooks and CI/CD events.

4.  **Storage Optimization**
    *   Implement `git gc` (Garbage Collection) scheduling.
    *   Implement object deduplication strategies for forks (alternates).

## Tech Stack
*   **Language:** Rust (Actix-Web or Axum)
*   **Libraries:** `git2` (libgit2), `russh` (SSH), `tokio` (Async runtime)

## Verification
*   User can run `git clone http://localhost:8081/repo/my-repo.git` successfully.
*   User can run `git push origin main` successfully.
*   Pushing to a protected branch is rejected with a meaningful error message.
