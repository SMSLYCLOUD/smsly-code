# Agent 6: CI/CD & Automation (SMSLY Actions)

**Mission:** Develop the integrated CI/CD system ("SMSLY Actions") to automatically build, test, and deploy code on every push or event.

## Goals

1.  **Workflow Parser (YAML)**
    *   Parse `.smsly/workflows/*.yml` files in the repository.
    *   Validate YAML syntax (Schema/Spec).
    *   Support `on:` triggers (push, pull_request, schedule, repository_dispatch).

2.  **Job Orchestration (Runner Service)**
    *   Develop a runner architecture (server-side coordination).
    *   Queue jobs based on `runs-on` (ubuntu-latest, self-hosted).
    *   Scale runners dynamically (Kubernetes/Docker).
    *   Implement "Self-Hosted Runners" capability (agent registration).

3.  **Secrets & Environments**
    *   Manage encrypted secrets (`/repos/{owner}/{repo}/settings/secrets`).
    *   Inject secrets into runner environment variables securely.
    *   Implement "Environment Protection" rules (e.g., approval for prod).

4.  **Artifacts & Logs**
    *   Stream build logs via WebSocket to the UI in real-time.
    *   Store build artifacts (binaries, test results) in MinIO/S3.
    *   Implement retention policies (e.g., 90 days).

5.  **Marketplace Actions**
    *   Support referencing reusable actions (`uses: actions/checkout@v4`).
    *   Create a registry for community actions (`marketplace.smsly.code`).

## Tech Stack
*   **Language:** Go (Orchestrator), Rust (Runner Agent - for speed/security)
*   **Database:** PostgreSQL (Job state), Redis (Queue)
*   **Storage:** MinIO (Artifacts, Logs)
*   **Compute:** Docker / Kubernetes (Job execution)

## Verification
*   Pushing code triggers a workflow defined in `.smsly/workflows/ci.yml`.
*   Job runs inside a Docker container.
*   Secrets are masked in logs (***).
*   User can download build artifacts from the UI.
