# Agent 2: API & Identity (Go)

**Mission:** Develop the authentication, authorization, and organization management layers. This agent ensures the correct users have the correct access to resources.

## Goals

1.  **Identity Provider (OIDC/SAML)**
    *   Implement standard OIDC/OAuth2 endpoints (`/api/auth/oauth2/authorize`, `/api/auth/oauth2/token`).
    *   Allow creating OAuth apps for integrations (Client ID/Secret management).
    *   (Optional) Support LDAP/SAML for enterprise environments.

2.  **Organization & Team Structure**
    *   Implement hierarchical models: Organizations → Teams → Members.
    *   RBAC Roles: Owner, Admin, Write, Read, None.
    *   Create endpoints:
        *   `POST /orgs`
        *   `POST /orgs/{org}/teams`
        *   `PUT /orgs/{org}/members`

3.  **SSH & GPG Key Management**
    *   Implement user settings to upload SSH public keys (`/user/settings/keys`).
    *   Implement GPG key storage for commit verification (`/user/settings/gpg`).
    *   Verify key uniqueness and format validity.

4.  **Fine-Grained Permissions (PATs)**
    *   Implement Personal Access Tokens (PATs) with specific scopes (e.g., `repo:read`, `workflow:write`).
    *   Enforce these scopes in middleware for all API requests.

5.  **Rate Limiting & Security**
    *   Implement global and per-user rate limits (Redis-backed).
    *   Implement audit logging for sensitive actions (e.g., changing repository visibility).

## Tech Stack
*   **Language:** Go (Fiber)
*   **Database:** PostgreSQL (GORM)
*   **Cache:** Redis
*   **Libraries:** `jwt-go`, `casbin` (RBAC policies)

## Verification
*   User can create an Organization and invite other users.
*   Team members inherit repository permissions correctly.
*   API requests fail appropriately with insufficient scopes.
*   SSH keys can be added and listed.
