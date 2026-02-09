# Phase 1 — Ready-to-Paste Jules Prompts

> Copy each prompt below into a separate Jules session.
> Deploy Cards 1.1, 2.1, 3.1, 4.1, 10.1 FIRST (zero dependencies).
> Then deploy the rest of Phase 1 after those merge.

---

## CARD 1.1 — Rust Git Engine: Repository Manager

```
REPO: https://github.com/SMSLYCLOUD/smsly-code
BRANCH: feature/sq1-card-1.1-repo-manager
PR TITLE: [SQ1] Card 1.1: Repository Manager

Read contracts/INTEGRATION_CONTRACTS.md first — it defines exact types and file paths.

## Task
Create the Rust Git engine crate at smsly-git/ with repository CRUD using libgit2.

## Files to Create
smsly-git/Cargo.toml
smsly-git/src/lib.rs
smsly-git/src/repo.rs
smsly-git/src/types.rs    — USE EXACT TYPES FROM INTEGRATION_CONTRACTS.md §3A
smsly-git/src/error.rs
smsly-git/src/config.rs
smsly-git/tests/repo_test.rs

## Cargo.toml Dependencies
[package]
name = "smsly-git"
version = "0.1.0"
edition = "2021"

[dependencies]
git2 = "0.19"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
chrono = { version = "0.4", features = ["serde"] }
uuid = { version = "1", features = ["v4", "serde"] }
thiserror = "2"
tracing = "0.1"
tempfile = "3"

[lib]
crate-type = ["cdylib", "rlib"]

## repo.rs Functions (EXACT signatures)
pub fn init_bare(base_path: &Path, owner: &str, name: &str) -> Result<RepoHandle, GitError>
pub fn open(base_path: &Path, owner: &str, name: &str) -> Result<RepoHandle, GitError>
pub fn delete(base_path: &Path, owner: &str, name: &str) -> Result<(), GitError>
pub fn exists(base_path: &Path, owner: &str, name: &str) -> bool
pub fn get_info(handle: &RepoHandle) -> Result<RepoInfo, GitError>
pub fn set_description(handle: &RepoHandle, desc: &str) -> Result<(), GitError>
pub fn set_default_branch(handle: &RepoHandle, branch: &str) -> Result<(), GitError>
pub fn fork(base_path: &Path, source: &RepoHandle, new_owner: &str, new_name: &str) -> Result<RepoHandle, GitError>

## error.rs
Use thiserror. Variants: NotFound, AlreadyExists, InvalidName, GitError(git2::Error), IoError, PermissionDenied

## types.rs
Copy ALL types EXACTLY from INTEGRATION_CONTRACTS.md section 3A. Do NOT change any field names or types.

## Tests (20+ required)
- init_bare creates valid bare repo
- open existing repo succeeds
- open non-existent repo returns NotFound
- delete removes repo directory
- delete non-existent returns NotFound
- exists returns true/false correctly
- get_info on empty repo
- get_info on repo with commits
- set_description persists
- set_default_branch persists
- fork creates independent copy
- init_bare with invalid name returns InvalidName
- concurrent init_bare (thread safety)
- Names with special chars rejected

## Acceptance Criteria
- [ ] cargo build succeeds with zero warnings
- [ ] cargo test passes all 20+ tests
- [ ] cargo clippy -- -D warnings passes
- [ ] All public items documented
- [ ] No unwrap() in library code
- [ ] Only creates files in smsly-git/
```

---

## CARD 2.1 — Go API Server Scaffold

```
REPO: https://github.com/SMSLYCLOUD/smsly-code
BRANCH: feature/sq2-card-2.1-api-scaffold
PR TITLE: [SQ2] Card 2.1: API Server Scaffold

Read contracts/INTEGRATION_CONTRACTS.md first.

## Task
Create the Go API server scaffold at smsly-code-api/ using Go Fiber v2.

## Files to Create
smsly-code-api/go.mod                          — module github.com/SMSLYCLOUD/smsly-code/smsly-code-api
smsly-code-api/cmd/server/main.go              — Entry point
smsly-code-api/internal/config/config.go       — Env config (uses envconfig)
smsly-code-api/internal/database/postgres.go   — pgxpool connection
smsly-code-api/internal/database/redis.go      — Redis connection
smsly-code-api/internal/middleware/auth.go      — JWT auth middleware (stub)
smsly-code-api/internal/middleware/cors.go
smsly-code-api/internal/middleware/ratelimit.go
smsly-code-api/internal/middleware/logger.go
smsly-code-api/internal/middleware/recovery.go
smsly-code-api/internal/middleware/requestid.go
smsly-code-api/internal/handlers/health.go     — GET /api/v1/health
smsly-code-api/router/router.go                — USE EXACT PATTERN FROM INTEGRATION_CONTRACTS.md §3D
smsly-code-api/pkg/response/response.go        — USE EXACT CODE FROM INTEGRATION_CONTRACTS.md §3B
smsly-code-api/pkg/validator/validator.go
smsly-code-api/migrations/000001_core_tables.sql

## Go Dependencies
github.com/gofiber/fiber/v2
github.com/jackc/pgx/v5
github.com/redis/go-redis/v9
github.com/golang-jwt/jwt/v5
github.com/rs/zerolog
github.com/google/uuid
github.com/kelseyhightower/envconfig

## Migration 000001_core_tables.sql
CREATE TABLE "user" (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(40) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255) DEFAULT '',
    avatar_url    VARCHAR(500) DEFAULT '',
    bio           TEXT DEFAULT '',
    location      VARCHAR(255) DEFAULT '',
    website       VARCHAR(255) DEFAULT '',
    is_admin      BOOLEAN DEFAULT FALSE,
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE TABLE repository (
    id             BIGSERIAL PRIMARY KEY,
    owner_id       BIGINT REFERENCES "user"(id) ON DELETE CASCADE,
    name           VARCHAR(100) NOT NULL,
    description    TEXT DEFAULT '',
    is_private     BOOLEAN DEFAULT FALSE,
    is_fork        BOOLEAN DEFAULT FALSE,
    fork_id        BIGINT REFERENCES repository(id),
    default_branch VARCHAR(255) DEFAULT 'main',
    stars          INT DEFAULT 0,
    forks          INT DEFAULT 0,
    size           BIGINT DEFAULT 0,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(owner_id, name)
);

CREATE TABLE organization (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(40) UNIQUE NOT NULL,
    full_name   VARCHAR(255) DEFAULT '',
    description TEXT DEFAULT '',
    avatar_url  VARCHAR(500) DEFAULT '',
    website     VARCHAR(255) DEFAULT '',
    location    VARCHAR(255) DEFAULT '',
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE team (
    id         BIGSERIAL PRIMARY KEY,
    org_id     BIGINT REFERENCES organization(id) ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    permission VARCHAR(20) DEFAULT 'read',
    UNIQUE(org_id, name)
);

## main.go
- Load config from env
- Connect to PostgreSQL (with retry)
- Connect to Redis (with retry)
- Setup Fiber app with middleware
- Register routes
- Graceful shutdown on SIGINT/SIGTERM
- Health endpoint returns {"data": {"status": "ok", "version": "0.1.0"}}

## router.go — Use Dependencies struct from INTEGRATION_CONTRACTS.md §3D
- For handlers that don't exist yet, use nil checks or comment out
- Only wire up HealthHandler initially
- Leave other routes as comments with the card number that will implement them

## Tests (10+)
- Health endpoint returns 200
- CORS middleware allows configured origins
- Rate limiter blocks after threshold
- Request ID middleware adds X-Request-ID header
- Recovery middleware catches panics

## Acceptance Criteria
- [ ] go build ./... succeeds
- [ ] go test ./... passes
- [ ] Server starts and /api/v1/health returns 200
- [ ] Only creates files in smsly-code-api/
```

---

## CARD 3.1 — Next.js Frontend Scaffold

```
REPO: https://github.com/SMSLYCLOUD/smsly-code
BRANCH: feature/sq3-card-3.1-frontend-scaffold
PR TITLE: [SQ3] Card 3.1: Frontend Scaffold + Design System

Read contracts/INTEGRATION_CONTRACTS.md first.

## Task
Initialize Next.js 14 project at smsly-code-web/ with design system.

## Commands
cd smsly-code-web
npx -y create-next-app@latest ./ --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --no-turbopack

## After scaffold, create:
src/lib/api.ts           — USE EXACT CODE FROM INTEGRATION_CONTRACTS.md §3E
src/lib/utils.ts         — cn() helper, formatDate, etc.
src/hooks/useAuth.ts     — Auth context (JWT management)
src/components/ui/Button.tsx
src/components/ui/Input.tsx
src/components/ui/Card.tsx
src/components/ui/Badge.tsx
src/components/ui/Avatar.tsx
src/components/ui/Skeleton.tsx
src/components/layout/Header.tsx
src/components/layout/Footer.tsx
src/app/layout.tsx       — Dark mode, Inter font, root layout
src/app/page.tsx         — Placeholder landing page with SMSLY Code branding

## Design System (in globals.css)
Brand colors: Primary #6366F1, Secondary #8B5CF6, Accent #06B6D4
Dark background: #0F172A, Secondary BG: #1E293B
Font: Inter from next/font/google

## Acceptance Criteria
- [ ] npm run dev starts without errors
- [ ] npm run build succeeds
- [ ] npm run lint passes
- [ ] Dark mode renders correctly
- [ ] All UI components render
- [ ] API client configured with types from contracts
- [ ] Only creates files in smsly-code-web/
```

---

## CARD 4.1 — MIP Core Library

```
REPO: https://github.com/SMSLYCLOUD/smsly-code
BRANCH: feature/sq4-card-4.1-mip-core
PR TITLE: [SQ4] Card 4.1: MIP Core Library

Read contracts/INTEGRATION_CONTRACTS.md first.

## Task
Create the MIP integrity stamp engine in Go.

## Files (ALL in smsly-code-api/internal/mip/)
stamp.go    — MIPStamp struct, CreateStamp()
merkle.go   — ComputeMerkleRoot() using SHA-256
verify.go   — VerifyStamp(), VerifyChain()
chain.go    — ChainVerification struct
crypto.go   — Ed25519 sign/verify (Go stdlib crypto/ed25519)
mip_test.go — 30+ tests

## Migration File
smsly-code-api/migrations/000004_mip_stamps.sql

CREATE TABLE mip_stamp (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id         BIGINT NOT NULL,
    commit_sha      VARCHAR(40) NOT NULL,
    merkle_root     VARCHAR(64) NOT NULL,
    tree_hash       VARCHAR(64) NOT NULL,
    author_id       BIGINT NOT NULL,
    parent_stamp_id UUID REFERENCES mip_stamp(id),
    signature       TEXT NOT NULL,
    verified        BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(repo_id, commit_sha)
);

CREATE INDEX idx_mip_stamp_repo ON mip_stamp(repo_id);
CREATE INDEX idx_mip_stamp_commit ON mip_stamp(commit_sha);

## Key Algorithms
- Merkle: sort files by path, SHA-256 each (path:hash), build binary tree
- Sign: Ed25519 over JSON(commit_sha+merkle_root+tree_hash+author_id+parent_stamp_id+timestamp)
- Chain: each stamp links to parent via parent_stamp_id

## Tests (30+): stamp creation, merkle determinism, signature verify, chain verify, tamper detection

## Acceptance Criteria
- [ ] go test ./internal/mip/... passes all 30+ tests
- [ ] Merkle root is deterministic
- [ ] Ed25519 signatures validate correctly
- [ ] Chain verification detects breaks
- [ ] Only creates files in smsly-code-api/internal/mip/ and migrations/
```

---

## CARD 10.1 — Docker Compose Setup

```
REPO: https://github.com/SMSLYCLOUD/smsly-code
BRANCH: feature/sq10-card-10.1-docker-compose
PR TITLE: [SQ10] Card 10.1: Docker Compose Development Setup

## Task
Create Docker Compose for local development with all infrastructure services.

## Files
docker-compose.dev.yml
docker/Dockerfile.api    — Multi-stage Go build
docker/Dockerfile.web    — Multi-stage Next.js build
docker/Dockerfile.git    — Rust build
docker/caddy/Caddyfile   — Reverse proxy config

## docker-compose.dev.yml Services
- postgres: PostgreSQL 16, port 5432, volume for data
- redis: Redis 7, port 6379
- minio: MinIO, ports 9000/9001, default credentials
- meilisearch: port 7700
- caddy: ports 80/443, reverse proxy to api+web

## Each service MUST have:
- health check
- restart: unless-stopped
- named volume for persistent data
- resource limits

## Acceptance Criteria
- [ ] docker compose -f docker-compose.dev.yml config validates
- [ ] All Dockerfiles have valid syntax
- [ ] Only creates files in docker/ and root docker-compose files
```
