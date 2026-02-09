# SMSLY Code — Contributing Guide

## Prerequisites

- **Rust** (latest stable) — for the Git engine
- **Go 1.22+** — for the API server
- **Node.js 20+** — for the frontend
- **Docker + Docker Compose** — for local services
- **PostgreSQL 16** — database (runs in Docker)
- **Redis 7** — cache/queue (runs in Docker)

## Quick Start

```bash
# Clone the repo
git clone https://github.com/SMSLYCLOUD/smsly-code.git
cd smsly-code

# Copy environment file
cp .env.example .env

# Start infrastructure services (Postgres, Redis, MinIO, Meilisearch)
make dev-up

# Build and run each component:
# Terminal 1: Rust Git engine (builds as shared library)
cd smsly-git && cargo build --release

# Terminal 2: Go API server
cd smsly-code-api && go run ./cmd/server

# Terminal 3: Next.js frontend
cd smsly-code-web && npm install && npm run dev
```

## Branch Naming

```
feature/sq{squadron}-card-{number}-{description}
```

Examples:
- `feature/sq1-card-1.1-repo-manager`
- `feature/sq4-card-4.1-mip-core`

## Commit Messages

```
[SQ{n}] Card {x.y}: {description}
```

Examples:
- `[SQ1] Card 1.1: Add repository CRUD operations`
- `[SQ4] Card 4.1: Implement MIP stamp creation and verification`

## PR Guidelines

- Title: `[SQ{n}] Card {x.y}: {description}`
- Body: Copy the card's acceptance criteria as a checklist
- All tests must pass
- No lint warnings

## Code Style

### Rust
- `cargo fmt` before committing
- `cargo clippy` must pass with zero warnings
- All public items must be documented
- No `unwrap()` in library code

### Go
- `gofmt` before committing
- `golangci-lint run` must pass
- All exported functions must be documented
- Use structured logging (zerolog)

### TypeScript
- ESLint + Prettier before committing
- No `any` types
- All components must have proper TypeScript props
