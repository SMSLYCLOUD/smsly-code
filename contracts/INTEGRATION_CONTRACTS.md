# SMSLY Code — Integration Contracts & Crash Prevention Guide

> **THIS IS THE MOST IMPORTANT DOCUMENT IN THE PROJECT.**
> Every Jules agent MUST read this before their card AND the Architecture Document.
> This prevents crashes when components integrate.

---

## THE PROBLEM THIS SOLVES

150 Jules agents build 150 components independently. Without contracts:
- Card 2.2 calls a function from Card 1.1 that has a different signature → **CRASH**
- Card 3.4 expects an API response format that Card 2.3 doesn't produce → **CRASH**
- Card 4.1 writes to a database table that Card 1.4 didn't create → **CRASH**
- Two cards create the same file → **MERGE CONFLICT**

This document defines the EXACT interfaces between all components so nothing crashes.

---

## RULE 1: MONOREPO STRUCTURE (IMMUTABLE)

```
smsly-code/                          ← Git root
├── README.md
├── LICENSE
├── Makefile                          ← Top-level build orchestration
├── docker-compose.yml               ← Production
├── docker-compose.dev.yml           ← Development
├── .env.example
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── build.yml
│       └── release.yml
│
├── smsly-git/                        ← SQUADRON 1 TERRITORY (Rust)
│   ├── Cargo.toml
│   ├── src/
│   │   ├── lib.rs
│   │   ├── repo.rs                   ← Card 1.1
│   │   ├── refs.rs                   ← Card 1.4
│   │   ├── objects.rs                ← Card 1.5
│   │   ├── walk.rs                   ← Card 1.5
│   │   ├── diff.rs                   ← Card 1.6
│   │   ├── error.rs                  ← Card 1.1
│   │   ├── types.rs                  ← Card 1.1
│   │   ├── config.rs                 ← Card 1.1
│   │   ├── transport/
│   │   │   ├── mod.rs
│   │   │   ├── http.rs               ← Card 1.2
│   │   │   ├── ssh.rs                ← Card 1.3
│   │   │   ├── pack.rs               ← Card 1.2
│   │   │   └── advertise.rs          ← Card 1.2
│   │   ├── hooks/
│   │   │   ├── mod.rs                ← Card 1.7
│   │   │   ├── engine.rs             ← Card 1.7
│   │   │   └── builtin.rs            ← Card 1.7
│   │   ├── lfs/
│   │   │   ├── mod.rs                ← Card 1.8
│   │   │   ├── batch.rs              ← Card 1.8
│   │   │   └── storage.rs            ← Card 1.8
│   │   ├── maintenance.rs            ← Card 1.9
│   │   └── ffi.rs                    ← Card 1.10
│   ├── smsly_git.h                   ← Card 1.10 (C header)
│   ├── tests/
│   └── benches/
│
├── smsly-code-api/                   ← SQUADRON 2 + 4-9,11-12 TERRITORY (Go)
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go               ← Card 2.1
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go             ← Card 2.1
│   │   ├── database/
│   │   │   ├── postgres.go           ← Card 2.1
│   │   │   └── redis.go              ← Card 2.1
│   │   ├── middleware/
│   │   │   ├── auth.go               ← Card 2.1
│   │   │   ├── cors.go               ← Card 2.1
│   │   │   ├── ratelimit.go          ← Card 2.1
│   │   │   ├── logger.go             ← Card 2.1
│   │   │   ├── recovery.go           ← Card 2.1
│   │   │   ├── requestid.go          ← Card 2.1
│   │   │   └── permission.go         ← Card 6.2
│   │   ├── models/
│   │   │   ├── user.go               ← Card 2.2
│   │   │   ├── repository.go         ← Card 2.3
│   │   │   ├── pull_request.go       ← Card 2.5
│   │   │   ├── issue.go              ← Card 2.6
│   │   │   ├── webhook.go            ← Card 2.7
│   │   │   ├── organization.go       ← Card 2.8
│   │   │   ├── release.go            ← Card 2.9
│   │   │   ├── ssh_key.go            ← Card 2.2
│   │   │   ├── token.go              ← Card 6.3
│   │   │   ├── mip_stamp.go          ← Card 4.1
│   │   │   ├── dip_certificate.go    ← Card 5.1
│   │   │   ├── trust_score.go        ← Card 4.1 (Squadron 4)
│   │   │   ├── notification.go       ← Card 9.3
│   │   │   ├── ai_review.go          ← Card 11.1
│   │   │   ├── workflow.go           ← Card 7.1
│   │   │   ├── audit_log.go          ← Card 12.1
│   │   │   └── deployment.go         ← Card 5.1
│   │   ├── handlers/
│   │   │   ├── health.go             ← Card 2.1
│   │   │   ├── auth.go               ← Card 2.2
│   │   │   ├── user.go               ← Card 2.2
│   │   │   ├── repository.go         ← Card 2.3
│   │   │   ├── content.go            ← Card 2.4
│   │   │   ├── pull_request.go       ← Card 2.5
│   │   │   ├── issue.go              ← Card 2.6
│   │   │   ├── webhook.go            ← Card 2.7
│   │   │   ├── organization.go       ← Card 2.8
│   │   │   ├── release.go            ← Card 2.9
│   │   │   ├── mip.go               ← Card 4.3
│   │   │   ├── dip.go               ← Card 5.3
│   │   │   ├── search.go            ← Card 8.1
│   │   │   ├── analytics.go         ← Card 8.2
│   │   │   ├── notification.go      ← Card 9.3
│   │   │   ├── ai_review.go         ← Card 11.1
│   │   │   └── audit.go             ← Card 12.1
│   │   ├── services/
│   │   │   ├── user_service.go       ← Card 2.2
│   │   │   ├── repo_service.go       ← Card 2.3
│   │   │   ├── git_service.go        ← Card 2.4 (calls FFI)
│   │   │   ├── pr_service.go         ← Card 2.5
│   │   │   ├── issue_service.go      ← Card 2.6
│   │   │   ├── webhook_service.go    ← Card 2.7
│   │   │   ├── org_service.go        ← Card 2.8
│   │   │   ├── mip_service.go        ← Card 4.1
│   │   │   ├── dip_service.go        ← Card 5.1
│   │   │   ├── trust_service.go      ← Card 4.1
│   │   │   ├── search_service.go     ← Card 8.1
│   │   │   ├── notification_service.go ← Card 9.2
│   │   │   ├── ai_review_service.go  ← Card 11.1
│   │   │   ├── cicd_service.go       ← Card 7.2
│   │   │   └── audit_service.go      ← Card 12.1
│   │   ├── mip/                      ← SQUADRON 4 TERRITORY
│   │   │   ├── stamp.go              ← Card 4.1
│   │   │   ├── merkle.go             ← Card 4.1
│   │   │   ├── verify.go             ← Card 4.1
│   │   │   ├── chain.go              ← Card 4.1
│   │   │   └── crypto.go             ← Card 4.1
│   │   ├── dip/                      ← SQUADRON 5 TERRITORY
│   │   │   ├── certificate.go        ← Card 5.1
│   │   │   ├── verify.go             ← Card 5.1
│   │   │   └── provenance.go         ← Card 5.1
│   │   ├── auth/                     ← SQUADRON 6 TERRITORY
│   │   │   ├── smsly_oauth.go        ← Card 6.1
│   │   │   ├── permission.go         ← Card 6.2
│   │   │   └── token.go              ← Card 6.3
│   │   ├── cicd/                     ← SQUADRON 7 TERRITORY
│   │   │   ├── parser.go             ← Card 7.1
│   │   │   ├── runner.go             ← Card 7.2
│   │   │   └── secrets.go            ← Card 7.4
│   │   ├── search/                   ← SQUADRON 8 TERRITORY
│   │   │   └── engine.go             ← Card 8.1
│   │   ├── events/                   ← SQUADRON 9 TERRITORY
│   │   │   ├── bus.go                ← Card 9.1
│   │   │   ├── handlers.go           ← Card 9.1
│   │   │   └── types.go              ← Card 9.1
│   │   ├── notifications/            ← SQUADRON 9 TERRITORY
│   │   │   ├── sms.go               ← Card 9.2
│   │   │   ├── voice.go             ← Card 9.2
│   │   │   ├── email.go             ← Card 9.5
│   │   │   └── websocket.go         ← Card 9.4
│   │   ├── ai/                       ← SQUADRON 11 TERRITORY
│   │   │   ├── review.go            ← Card 11.1
│   │   │   ├── summary.go           ← Card 11.2
│   │   │   └── scanner.go           ← Card 11.3
│   │   └── enterprise/               ← SQUADRON 12 TERRITORY
│   │       ├── audit.go              ← Card 12.1
│   │       ├── saml.go              ← Card 12.2
│   │       ├── ip_allowlist.go      ← Card 12.3
│   │       └── rbac.go              ← Card 12.4
│   ├── router/
│   │   └── router.go                ← Card 2.1 (ALL routes registered here)
│   ├── pkg/
│   │   ├── response/
│   │   │   └── response.go          ← Card 2.1
│   │   └── validator/
│   │       └── validator.go          ← Card 2.1
│   └── migrations/
│       ├── 000001_core_tables.sql    ← Card 2.1 (users, repos, orgs)
│       ├── 000002_pr_issues.sql      ← Card 2.5 (PRs, reviews, issues)
│       ├── 000003_smsly_tables.sql   ← Card 1.4 (MIP, DIP, trust, ai_reviews)
│       └── 000004_enterprise.sql     ← Card 12.1 (audit_log, roles)
│
├── smsly-code-web/                   ← SQUADRON 3 TERRITORY (Next.js)
│   ├── package.json
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── src/
│       ├── app/                      ← All Card 3.x pages
│       ├── components/               ← All Card 3.x components
│       ├── lib/                      ← Card 3.1
│       ├── hooks/                    ← Card 3.1
│       └── styles/                   ← Card 3.1
│
├── docker/                           ← SQUADRON 10 TERRITORY
│   ├── Dockerfile.api
│   ├── Dockerfile.web
│   ├── Dockerfile.git
│   └── caddy/
│       └── Caddyfile
│
├── deploy/                           ← SQUADRON 10 TERRITORY
│   ├── k8s/
│   └── scripts/
│
├── docs/                             ← SQUADRON 10 TERRITORY
│   ├── book.toml
│   └── src/
│
└── contracts/                        ← THIS FILE LIVES HERE
    └── INTEGRATION_CONTRACTS.md
```

**RULE: Each card ONLY creates/modifies files in their assigned location above. NEVER touch another squadron's files.**

---

## RULE 2: DEPENDENCY CHAIN (Build Order)

Cards MUST be deployed in this order. A card with dependencies MUST NOT start until its dependencies are merged.

### PHASE 1 — FOUNDATIONS (No dependencies, all deploy simultaneously)

```
WEEK 1, DAY 1-2: These cards have ZERO dependencies — deploy all at once:

Card 1.1  (Repo Manager)        → Creates: smsly-git/src/{lib,repo,error,types,config}.rs
Card 2.1  (API Scaffold)        → Creates: smsly-code-api/ entire scaffold
Card 3.1  (Frontend Scaffold)   → Creates: smsly-code-web/ entire scaffold
Card 4.1  (MIP Core Library)    → Creates: smsly-code-api/internal/mip/
Card 5.1  (DIP Certificate)     → Creates: smsly-code-api/internal/dip/
Card 6.1  (SMSLY Auth)          → Creates: smsly-code-api/internal/auth/smsly_oauth.go
Card 7.1  (Workflow Parser)     → Creates: smsly-code-api/internal/cicd/parser.go
Card 8.1  (Search Engine)       → Creates: smsly-code-api/internal/search/engine.go
Card 9.1  (Event Bus)           → Creates: smsly-code-api/internal/events/
Card 10.1 (Docker Compose)      → Creates: docker-compose.yml, .env.example, Makefile
Card 12.1 (Audit Log)           → Creates: smsly-code-api/internal/enterprise/audit.go

WHY NO CRASHES: Each card creates its OWN directory/files. No overlaps.
```

### PHASE 2 — DEPENDS ON PHASE 1

```
WEEK 1, DAY 3-5: These depend on specific Phase 1 cards:

Card 1.2  (HTTP Transport)   ← Depends on: Card 1.1 (needs RepoHandle type)
Card 1.3  (SSH Transport)    ← Depends on: Card 1.1 (needs RepoHandle type)
Card 1.4  (Ref Management)   ← Depends on: Card 1.1 (needs RepoHandle type)
Card 1.5  (Objects/History)  ← Depends on: Card 1.1 (needs RepoHandle type)
Card 2.2  (User/Auth API)    ← Depends on: Card 2.1 (needs scaffold)
Card 2.3  (Repository API)   ← Depends on: Card 2.1 (needs scaffold)
Card 3.2  (Landing Page)     ← Depends on: Card 3.1 (needs scaffold)
Card 3.3  (Auth Pages)       ← Depends on: Card 3.1 (needs scaffold)
Card 6.2  (Permissions)      ← Depends on: Card 2.1 (needs middleware pattern)
Card 6.3  (API Tokens)       ← Depends on: Card 2.2 (needs user model)
Card 9.2  (SMS/Voice)        ← Depends on: Card 9.1 (needs event bus)
Card 10.5 (CI Pipeline)      ← Depends on: Card 10.1 (needs Docker files)

WHY NO CRASHES: Each card reads from Phase 1 outputs but creates NEW files.
```

### PHASE 3 — DEPENDS ON PHASE 2

```
WEEK 2: These depend on Phase 2 cards:

Card 1.6  (Diff Engine)     ← Depends on: Card 1.5 (needs objects.rs)
Card 1.7  (Hooks Engine)    ← Depends on: Card 1.1 (needs RepoHandle)
Card 1.8  (LFS Server)      ← Depends on: Card 1.2 (needs HTTP transport)
Card 1.10 (FFI Layer)       ← Depends on: Cards 1.1-1.6 (wraps all functions)
Card 2.4  (File Browsing)   ← Depends on: Card 2.3 + Card 1.10 (needs FFI)
Card 2.5  (Pull Request)    ← Depends on: Card 2.3 (needs repo model)
Card 2.6  (Issues)          ← Depends on: Card 2.3 (needs repo model)
Card 2.7  (Webhooks)        ← Depends on: Card 2.3 + Card 9.1 (needs event bus)
Card 2.8  (Organizations)   ← Depends on: Card 2.2 (needs user model)
Card 2.9  (Releases)        ← Depends on: Card 2.3 (needs repo model)
Card 3.4  (Repo Page)       ← Depends on: Card 3.1 + Card 2.4 (needs API)
Card 3.5  (File Viewer)     ← Depends on: Card 3.1 + Card 2.4 (needs API)
Card 4.2  (MIP Git Hooks)   ← Depends on: Card 4.1 + Card 1.7 (needs hooks engine)
Card 5.2  (DIP CI)          ← Depends on: Card 5.1 + Card 7.2 (needs CI runner)
Card 7.2  (Job Runner)      ← Depends on: Card 7.1 (needs parser)
Card 9.3  (Notif Prefs)     ← Depends on: Card 9.2 + Card 2.2 (needs user model)
```

### PHASE 4 — DEPENDS ON PHASE 3

```
WEEK 2-3: Full integration cards:

Card 1.9  (Maintenance)      ← Depends on: Card 1.1
Card 3.6  (Diff Viewer)      ← Depends on: Card 3.1 + Card 2.5 (needs PR API)
Card 3.7  (PR Page)          ← Depends on: Card 3.6 + Card 2.5
Card 3.8  (Issues Page)      ← Depends on: Card 3.1 + Card 2.6
Card 3.9  (Profile/Settings) ← Depends on: Card 3.1 + Card 2.2
Card 3.10 (Integrity Dash)   ← Depends on: Card 3.1 + Card 4.3 + Card 5.3
Card 4.3  (MIP API)          ← Depends on: Card 4.1 + Card 2.1
Card 5.3  (DIP API)          ← Depends on: Card 5.1 + Card 2.1
Card 7.3  (Build Status API) ← Depends on: Card 7.2 + Card 2.1
Card 8.2  (Analytics)        ← Depends on: Card 2.3
Card 9.4  (WebSocket)        ← Depends on: Card 9.1 + Card 2.1
Card 9.5  (Email Templates)  ← Depends on: Card 9.2
Card 11.1 (AI Code Review)   ← Depends on: Card 2.5 (needs PR model)
Card 11.2 (PR Summary)       ← Depends on: Card 2.5
Card 11.3 (Vuln Scanner)     ← Depends on: Card 2.3
Card 12.2 (SAML SSO)         ← Depends on: Card 2.2 + Card 6.1
Card 12.3 (IP Allowlists)    ← Depends on: Card 2.8
Card 12.4 (Custom Roles)     ← Depends on: Card 6.2
```

### PHASE 5 — POLISH (Week 4+)

All remaining cards: testing, docs, frontend polish, UX improvements.

---

## RULE 3: SHARED TYPE CONTRACTS (Exact signatures — DO NOT DEVIATE)

### 3A. Rust Types (smsly-git/src/types.rs) — Created by Card 1.1

Every Squadron 1 card MUST use these EXACT types:

```rust
// ============= THESE ARE CANONICAL — DO NOT CHANGE =============

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::path::PathBuf;

/// Handle to an open Git repository
pub struct RepoHandle {
    pub path: PathBuf,
    pub owner: String,
    pub name: String,
    pub(crate) inner: git2::Repository,
}

/// Git commit information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CommitInfo {
    pub sha: String,
    pub short_sha: String,
    pub message: String,
    pub body: Option<String>,
    pub author: Signature,
    pub committer: Signature,
    pub parents: Vec<String>,
    pub tree_sha: String,
    pub timestamp: DateTime<Utc>,
    pub is_merge: bool,
}

/// Git author/committer signature
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Signature {
    pub name: String,
    pub email: String,
    pub timestamp: DateTime<Utc>,
}

/// Branch information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BranchInfo {
    pub name: String,
    pub sha: String,
    pub is_default: bool,
    pub is_protected: bool,
}

/// Tag information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TagInfo {
    pub name: String,
    pub sha: String,
    pub is_annotated: bool,
    pub message: Option<String>,
    pub tagger: Option<Signature>,
    pub target_sha: String,
}

/// File tree entry
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TreeEntry {
    pub name: String,
    pub path: String,
    pub sha: String,
    pub entry_type: EntryType,
    pub size: Option<u64>,
    pub mode: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EntryType {
    Blob,
    Tree,
    Submodule,
    Symlink,
}

/// Diff result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffResult {
    pub files: Vec<DiffFile>,
    pub stats: DiffStats,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffStats {
    pub additions: usize,
    pub deletions: usize,
    pub files_changed: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffFile {
    pub old_path: Option<String>,
    pub new_path: Option<String>,
    pub status: DiffStatus,
    pub hunks: Vec<DiffHunk>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum DiffStatus {
    Added,
    Modified,
    Deleted,
    Renamed,
    Copied,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffHunk {
    pub old_start: u32,
    pub old_lines: u32,
    pub new_start: u32,
    pub new_lines: u32,
    pub header: String,
    pub lines: Vec<DiffLine>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DiffLine {
    pub line_type: LineType,
    pub content: String,
    pub old_lineno: Option<u32>,
    pub new_lineno: Option<u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum LineType {
    Add,
    Delete,
    Context,
}

/// Repository information
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RepoInfo {
    pub owner: String,
    pub name: String,
    pub default_branch: String,
    pub size_bytes: u64,
    pub branch_count: usize,
    pub tag_count: usize,
    pub last_commit_at: Option<DateTime<Utc>>,
    pub is_empty: bool,
}

/// File content
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FileContent {
    pub content: Vec<u8>,
    pub size: u64,
    pub encoding: String,
    pub is_binary: bool,
    pub mime_type: String,
    pub sha: String,
}

/// Blame line
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlameLine {
    pub line_number: usize,
    pub content: String,
    pub commit_sha: String,
    pub author: Signature,
    pub date: DateTime<Utc>,
}

/// Hook types
pub enum HookType {
    PreReceive,
    Update,
    PostReceive,
}

/// Ref update (used in hooks)
#[derive(Debug, Clone)]
pub struct RefUpdate {
    pub ref_name: String,
    pub old_sha: String,
    pub new_sha: String,
}
```

### 3B. Go API Response Format — Created by Card 2.1

Every Go handler MUST use these EXACT response helpers:

```go
package response

import "github.com/gofiber/fiber/v2"

// Standard success response
type APIResponse struct {
    Data  interface{} `json:"data,omitempty"`
    Meta  *Meta       `json:"meta,omitempty"`
    Error *APIError   `json:"error,omitempty"`
}

type Meta struct {
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
    Total   int `json:"total"`
}

type APIError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
    return c.JSON(APIResponse{Data: data})
}

func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *Meta) error {
    return c.JSON(APIResponse{Data: data, Meta: meta})
}

func Error(c *fiber.Ctx, status int, code, message string) error {
    return c.Status(status).JSON(APIResponse{
        Error: &APIError{Code: code, Message: message},
    })
}

func NotFound(c *fiber.Ctx, resource string) error {
    return Error(c, 404, "NOT_FOUND", resource+" not found")
}

func Unauthorized(c *fiber.Ctx) error {
    return Error(c, 401, "UNAUTHORIZED", "Authentication required")
}

func Forbidden(c *fiber.Ctx) error {
    return Error(c, 403, "FORBIDDEN", "Insufficient permissions")
}

func BadRequest(c *fiber.Ctx, message string) error {
    return Error(c, 400, "BAD_REQUEST", message)
}

func InternalError(c *fiber.Ctx, err error) error {
    // Log the actual error, return generic message
    log.Error().Err(err).Msg("Internal server error")
    return Error(c, 500, "INTERNAL_ERROR", "An internal error occurred")
}
```

### 3C. Go Database Access Pattern — ALL Go cards use this

```go
// EVERY Go card that accesses the database MUST:
// 1. Accept *pgxpool.Pool as dependency injection (NOT global variable)
// 2. Use context.Context for all queries
// 3. Use pgx v5 syntax
// 4. Return typed errors

// Example service pattern (EVERY service MUST follow this):
type UserService struct {
    db    *pgxpool.Pool
    redis *redis.Client
}

func NewUserService(db *pgxpool.Pool, redis *redis.Client) *UserService {
    return &UserService{db: db, redis: redis}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := s.db.QueryRow(ctx,
        `SELECT id, username, email, full_name, avatar_url, bio, 
                location, website, is_admin, created_at, last_login_at
         FROM "user" WHERE id = $1`, id,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.FullName,
           &user.AvatarURL, &user.Bio, &user.Location, &user.Website,
           &user.IsAdmin, &user.CreatedAt, &user.LastLoginAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("get user by id %d: %w", id, err)
    }
    return &user, nil
}
```

### 3D. Go Router Registration Pattern — Card 2.1 creates, ALL handlers register here

```go
// router/router.go — ONLY Card 2.1 creates this file
// ALL other cards add routes via registration functions

package router

func SetupRoutes(app *fiber.App, deps *Dependencies) {
    api := app.Group("/api/v1")

    // Health — Card 2.1
    api.Get("/health", deps.HealthHandler.Health)

    // Auth — Card 2.2
    auth := api.Group("/auth")
    auth.Post("/register", deps.AuthHandler.Register)
    auth.Post("/login", deps.AuthHandler.Login)
    auth.Post("/refresh", deps.AuthHandler.Refresh)
    auth.Post("/logout", deps.AuthHandler.Logout)

    // User — Card 2.2
    user := api.Group("/user", deps.AuthMiddleware)
    user.Get("/", deps.UserHandler.GetCurrent)
    user.Patch("/", deps.UserHandler.Update)
    user.Get("/repos", deps.RepoHandler.ListUserRepos)
    user.Get("/keys", deps.UserHandler.ListKeys)
    user.Post("/keys", deps.UserHandler.AddKey)
    user.Get("/tokens", deps.TokenHandler.List)
    user.Post("/tokens", deps.TokenHandler.Create)

    // Repos — Card 2.3+
    repos := api.Group("/repos")
    repos.Post("/", deps.AuthMiddleware, deps.RepoHandler.Create)
    repos.Get("/search", deps.RepoHandler.Search)

    // Repo routes — Card 2.3, 2.4, 2.5, 2.6, 2.7
    repo := repos.Group("/:owner/:repo", deps.RepoMiddleware)
    repo.Get("/", deps.RepoHandler.Get)
    repo.Patch("/", deps.RepoPermission("admin"), deps.RepoHandler.Update)
    repo.Delete("/", deps.RepoPermission("admin"), deps.RepoHandler.Delete)

    // Contents — Card 2.4
    repo.Get("/contents/*", deps.ContentHandler.GetContents)
    repo.Get("/commits", deps.ContentHandler.ListCommits)
    repo.Get("/commits/:sha", deps.ContentHandler.GetCommit)
    repo.Get("/branches", deps.ContentHandler.ListBranches)
    repo.Get("/tags", deps.ContentHandler.ListTags)
    repo.Get("/compare/:base...:head", deps.ContentHandler.Compare)

    // PRs — Card 2.5
    repo.Post("/pulls", deps.AuthMiddleware, deps.PRHandler.Create)
    repo.Get("/pulls", deps.PRHandler.List)
    repo.Get("/pulls/:number", deps.PRHandler.Get)
    repo.Post("/pulls/:number/merge", deps.RepoPermission("write"), deps.PRHandler.Merge)

    // Issues — Card 2.6
    repo.Post("/issues", deps.AuthMiddleware, deps.IssueHandler.Create)
    repo.Get("/issues", deps.IssueHandler.List)
    repo.Get("/issues/:number", deps.IssueHandler.Get)

    // MIP — Card 4.3
    repo.Get("/stamps", deps.MIPHandler.ListStamps)
    repo.Get("/stamps/:sha", deps.MIPHandler.GetStamp)
    repo.Post("/stamps/:sha/verify", deps.MIPHandler.VerifyStamp)
    repo.Get("/stamps/chain", deps.MIPHandler.GetChain)

    // DIP — Card 5.3
    repo.Get("/certificates", deps.DIPHandler.ListCerts)
    repo.Get("/certificates/:id", deps.DIPHandler.GetCert)
    repo.Get("/provenance/:sha", deps.DIPHandler.GetProvenance)

    // Webhooks — Card 2.7
    repo.Post("/hooks", deps.RepoPermission("admin"), deps.WebhookHandler.Create)
    repo.Get("/hooks", deps.RepoPermission("admin"), deps.WebhookHandler.List)

    // Trust Score — Card 4.3
    repo.Get("/trust-score", deps.TrustHandler.GetScore)

    // AI Review — Card 11.1
    repo.Get("/pulls/:number/ai-review", deps.AIReviewHandler.Get)
    repo.Post("/pulls/:number/ai-review", deps.AIReviewHandler.Trigger)

    // Analytics — Card 8.2
    repo.Get("/analytics/commits", deps.AnalyticsHandler.CommitFrequency)
    repo.Get("/analytics/languages", deps.AnalyticsHandler.Languages)

    // Deploy — Card 5.3
    deploy := api.Group("/deploy")
    deploy.Post("/verify", deps.DIPHandler.VerifyDeploy)

    // Search — Card 8.1
    search := api.Group("/search")
    search.Get("/code", deps.SearchHandler.SearchCode)
    search.Get("/repos", deps.SearchHandler.SearchRepos)

    // Notifications — Card 9.3
    notifs := api.Group("/user/notifications", deps.AuthMiddleware)
    notifs.Get("/", deps.NotificationHandler.List)
    notifs.Get("/preferences", deps.NotificationHandler.GetPreferences)
    notifs.Put("/preferences", deps.NotificationHandler.UpdatePreferences)

    // Organizations — Card 2.8
    orgs := api.Group("/orgs")
    orgs.Post("/", deps.AuthMiddleware, deps.OrgHandler.Create)
    orgs.Get("/:org", deps.OrgHandler.Get)

    // Admin — Card 12.1
    admin := api.Group("/admin", deps.AdminMiddleware)
    admin.Get("/audit-log", deps.AuditHandler.List)

    // WebSocket — Card 9.4
    app.Get("/ws", deps.WebSocketHandler.Handle)

    // Git Smart HTTP — Card 1.2 (via Go reverse proxy to Rust)
    app.Get("/:owner/:repo.git/info/refs", deps.GitHTTPHandler.InfoRefs)
    app.Post("/:owner/:repo.git/git-upload-pack", deps.GitHTTPHandler.UploadPack)
    app.Post("/:owner/:repo.git/git-receive-pack", deps.GitHTTPHandler.ReceivePack)
}

// Dependencies struct — EVERY handler registers here
type Dependencies struct {
    AuthMiddleware   fiber.Handler
    AdminMiddleware  fiber.Handler
    RepoMiddleware   fiber.Handler
    RepoPermission   func(level string) fiber.Handler

    HealthHandler      *handlers.HealthHandler
    AuthHandler        *handlers.AuthHandler
    UserHandler        *handlers.UserHandler
    RepoHandler        *handlers.RepoHandler
    ContentHandler     *handlers.ContentHandler
    PRHandler          *handlers.PRHandler
    IssueHandler       *handlers.IssueHandler
    WebhookHandler     *handlers.WebhookHandler
    OrgHandler         *handlers.OrgHandler
    ReleaseHandler     *handlers.ReleaseHandler
    MIPHandler         *handlers.MIPHandler
    DIPHandler         *handlers.DIPHandler
    TrustHandler       *handlers.TrustHandler
    SearchHandler      *handlers.SearchHandler
    AnalyticsHandler   *handlers.AnalyticsHandler
    NotificationHandler *handlers.NotificationHandler
    AIReviewHandler    *handlers.AIReviewHandler
    AuditHandler       *handlers.AuditHandler
    WebSocketHandler   *handlers.WebSocketHandler
    GitHTTPHandler     *handlers.GitHTTPHandler
    TokenHandler       *handlers.TokenHandler
}
```

### 3E. Frontend API Client — Card 3.1 creates, ALL frontend cards use

```typescript
// lib/api.ts — EVERY frontend card imports from here
// Card 3.1 creates this. Other cards ONLY add types.

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface APIResponse<T> {
  data?: T;
  meta?: { page: number; per_page: number; total: number };
  error?: { code: string; message: string; details?: any };
}

class APIClient {
  private token: string | null = null;

  setToken(token: string) { this.token = token; }
  clearToken() { this.token = null; }

  private async request<T>(method: string, path: string, body?: any): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (this.token) headers['Authorization'] = `Bearer ${this.token}`;

    const res = await fetch(`${BASE_URL}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    const json: APIResponse<T> = await res.json();
    if (json.error) throw new APIError(json.error.code, json.error.message, res.status);
    return json.data as T;
  }

  async get<T>(path: string): Promise<T> { return this.request('GET', path); }
  async post<T>(path: string, body: any): Promise<T> { return this.request('POST', path, body); }
  async patch<T>(path: string, body: any): Promise<T> { return this.request('PATCH', path, body); }
  async delete(path: string): Promise<void> { await this.request('DELETE', path); }
}

export const api = new APIClient();

export class APIError extends Error {
  constructor(public code: string, message: string, public status: number) {
    super(message);
  }
}

// === TYPE DEFINITIONS (each card adds its types here) ===

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  avatar_url: string;
  bio: string;
  location: string;
  website: string;
  is_admin: boolean;
  created_at: string;
}

export interface Repository {
  id: number;
  owner: User;
  name: string;
  full_name: string;
  description: string;
  is_private: boolean;
  is_fork: boolean;
  default_branch: string;
  stars: number;
  forks: number;
  open_issues: number;
  open_prs: number;
  size: number;
  language: string;
  trust_score: number;
  clone_url: string;
  ssh_url: string;
  created_at: string;
  updated_at: string;
}

export interface CommitInfo {
  sha: string;
  short_sha: string;
  message: string;
  author: { name: string; email: string; timestamp: string };
  committer: { name: string; email: string; timestamp: string };
  parents: string[];
  tree_sha: string;
  is_merge: boolean;
}

export interface TreeEntry {
  name: string;
  path: string;
  sha: string;
  entry_type: 'Blob' | 'Tree' | 'Submodule' | 'Symlink';
  size?: number;
  mode: number;
}

export interface PullRequest {
  id: number;
  number: number;
  title: string;
  body: string;
  state: 'open' | 'closed' | 'merged';
  author: User;
  base_branch: string;
  head_branch: string;
  additions: number;
  deletions: number;
  changed_files: number;
  labels: Label[];
  created_at: string;
  updated_at: string;
}

export interface Issue {
  id: number;
  number: number;
  title: string;
  body: string;
  state: 'open' | 'closed';
  author: User;
  assignees: User[];
  labels: Label[];
  created_at: string;
}

export interface Label {
  id: number;
  name: string;
  color: string;
  description: string;
}

export interface MIPStamp {
  id: string;
  commit_sha: string;
  merkle_root: string;
  author_id: number;
  verified: boolean;
  created_at: string;
}

export interface DIPCertificate {
  id: string;
  source_commit_sha: string;
  artifact_hash: string;
  artifact_type: string;
  environment: string;
  deployed_at?: string;
  created_at: string;
}

export interface TrustScore {
  score: number;
  mip_coverage: number;
  dip_coverage: number;
  test_coverage: number;
  review_coverage: number;
}
```

---

## RULE 4: DATABASE MIGRATION ORDER

Migrations MUST be numbered sequentially. Each card that needs DB changes creates a NEW migration file:

```
migrations/
├── 000001_core_tables.sql      ← Card 2.1 (user, repository, organization, team)
├── 000002_pr_issues.sql        ← Card 2.5 (pull_request, review, issue, comment, label, milestone)
├── 000003_webhooks.sql         ← Card 2.7 (webhook, webhook_delivery)
├── 000004_mip_stamps.sql       ← Card 4.1 (mip_stamps)
├── 000005_dip_certs.sql        ← Card 5.1 (dip_certificates)
├── 000006_trust_score.sql      ← Card 4.1 (user_trust)
├── 000007_notifications.sql    ← Card 9.3 (notification_preferences, notifications)
├── 000008_ai_reviews.sql       ← Card 11.1 (ai_reviews)
├── 000009_deployments.sql      ← Card 5.1 (deployments)
├── 000010_workflows.sql        ← Card 7.1 (workflow_run, workflow_job, workflow_step)
├── 000011_audit_log.sql        ← Card 12.1 (audit_log)
├── 000012_releases.sql         ← Card 2.9 (release, release_asset)
├── 000013_enterprise.sql       ← Card 12.4 (custom_role, role_permission, ip_allowlist)
└── 000014_ssh_keys.sql         ← Card 2.2 (ssh_key, api_token)
```

**RULE: Each card creates ONLY its own migration file. NEVER modify another card's migration.**

---

## RULE 5: Go Module Path

ALL Go imports MUST use:
```
github.com/SMSLYCLOUD/smsly-code/smsly-code-api
```

Example:
```go
import (
    "github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/models"
    "github.com/SMSLYCLOUD/smsly-code/smsly-code-api/internal/services"
    "github.com/SMSLYCLOUD/smsly-code/smsly-code-api/pkg/response"
)
```

---

## RULE 6: WHAT TO DO WHEN A DEPENDENCY ISN'T MERGED YET

If your card depends on another card's code that isn't merged yet:

1. **Create stub interfaces** for the dependency
2. **Use dependency injection** — accept interfaces, not concrete types
3. **Write tests with mocks** of the dependency
4. Your card MUST compile and pass tests INDEPENDENTLY

Example: Card 2.4 (File Browsing) needs Card 1.10 (FFI) which may not exist yet:

```go
// Define an interface for what you need:
type GitEngine interface {
    GetTree(owner, repo, branch, path string) ([]TreeEntry, error)
    GetBlob(owner, repo, sha string) ([]byte, error)
    ListCommits(owner, repo, branch string, page, perPage int) ([]CommitInfo, int, error)
}

// Your handler accepts the interface:
type ContentHandler struct {
    git GitEngine
}

// In tests, use a mock:
type MockGitEngine struct{}
func (m *MockGitEngine) GetTree(...) ([]TreeEntry, error) {
    return []TreeEntry{{Name: "README.md", EntryType: "Blob"}}, nil
}
```

**This way your card compiles and tests pass even if the dependency isn't merged.**

---

## RULE 7: BRANCH NAMING

Every Jules card creates a branch named:
```
feature/sq{squadron}-card-{number}-{short-description}
```

Examples:
```
feature/sq1-card-1.1-repo-manager
feature/sq2-card-2.1-api-scaffold
feature/sq3-card-3.1-frontend-scaffold
feature/sq4-card-4.1-mip-core
```

PR title: `[SQ{n}] Card {x.y}: {description}`
Example: `[SQ1] Card 1.1: Repository Manager`

---

## RULE 8: ENVIRONMENT VARIABLES

ALL environment variables MUST be prefixed. Cards MUST NOT invent new prefixes:

```
SMSLY_CODE_*     — App configuration
DATABASE_*       — PostgreSQL
REDIS_*          — Redis
MINIO_*          — Object storage
SMSLY_IDENTITY_* — SMSLY SSO
SMSLY_GATEWAY_*  — SMS/Voice notifications
GEMINI_*         — AI features
MEILISEARCH_*    — Search
```

---

## RULE 9: TESTING REQUIREMENTS

Every card MUST:
1. Include a `*_test.go` or `*_test.rs` or `*.test.tsx` file
2. Tests MUST pass independently (`go test ./internal/mip/...` must work)
3. Tests MUST NOT depend on external services (use mocks/stubs)
4. Tests MUST NOT require other cards to be merged
5. Minimum test counts are specified in each card

---

## RULE 10: CONFLICT RESOLUTION

If two cards both need to modify the same file (e.g., router.go):

1. **Card 2.1 creates router.go** with stub handler registration
2. **Other cards do NOT modify router.go directly**
3. Instead, each card creates its handler file (e.g., `handlers/mip.go`)
4. Integration PR (done by YOU, not Jules) wires it all together

**Alternatively:** Card 2.1 creates a registration mechanism:
```go
// Each handler package exports a RegisterRoutes function:
func RegisterMIPRoutes(router fiber.Router, handler *MIPHandler) {
    router.Get("/stamps", handler.ListStamps)
    router.Get("/stamps/:sha", handler.GetStamp)
}
```
Then Card 2.1's router.go calls all registration functions.
