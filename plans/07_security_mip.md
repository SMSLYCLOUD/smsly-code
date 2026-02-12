# Agent 7: Security & Compliance (MIP/DIP - The Core USP)

**Mission:** Build the unique selling point of SMSLY Code: Cryptographically verified integrity and deployment provenance, far beyond standard Git hosting.

## Goals

1.  **MIP: Manual Integrity Protocol (Signatures)**
    *   **Goal:** Ensure *every* commit is signed and linked to a verified identity.
    *   Implement strict GPG/SSH signature enforcement (`git config commit.gpgsign true`).
    *   Implement Sigstore/Rekor integration (keyless signing).
    *   Display "Verified" badges prominently in the UI.

2.  **DIP: Deployment Integrity Protocol (Provenance)**
    *   **Goal:** Trace every running artifact back to its source commit.
    *   Implement attestation generation during CI/CD builds (SLSA Level 3).
    *   Store attestations in a transparency log (Rekor).
    *   Create a "Deployments" view showing commit → build → artifact → environment.

3.  **Dependency Graph & Supply Chain**
    *   Parse `package.json`, `go.mod`, `Cargo.toml`, etc.
    *   Build a dependency graph for repositories.
    *   Integrate with vulnerability databases (OSV, CVE).
    *   Alert users on vulnerable dependencies (Dependabot-like).

4.  **Secret Scanning**
    *   Scan commits on push for regex patterns of known secrets (AWS keys, tokens).
    *   Reject pushes containing secrets (pre-receive hook) or alert admins.

5.  **Audit Logs**
    *   Log every meaningful action (repo created, member added, branch deleted).
    *   Provide an exportable audit log for compliance (`/orgs/{org}/settings/audit-log`).

## Tech Stack
*   **Language:** Rust (Scanning/Verification), Go (API/Logging)
*   **Libraries:** `sigstore-go`, `trivy` (scanning), `yara` (pattern matching)
*   **Database:** ClickHouse or Elastic (High-volume logs)

## Verification
*   Commits without a valid signature are flagged "Unverified" (red).
*   A deployment pipeline generates a signed attestation file.
*   Pushing a fake AWS key is blocked by the server.
*   Dependency graph correctly identifies `react` version in a JS project.
