# SMSLY Code

> **Where every commit is a promise.**

[![Build Status](https://github.com/SMSLYCLOUD/smsly-code/actions/workflows/ci.yml/badge.svg)](https://github.com/SMSLYCLOUD/smsly-code/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**SMSLY Code** is an integrity-verified Git hosting platform. Every commit is cryptographically stamped (MIP), every deployment is traced back to its exact source (DIP), and every repository has a real-time trust score.

## What Makes SMSLY Code Different

| Feature | GitHub | GitLab | **SMSLY Code** |
|---------|--------|--------|----------------|
| Git hosting | ✅ | ✅ | ✅ |
| CI/CD | ✅ | ✅ | ✅ |
| AI Code Review | Copilot | Duo | ✅ Gemini-powered |
| **MIP Integrity Stamps** | ❌ | ❌ | ✅ Every commit |
| **DIP Deploy Provenance** | ❌ | ❌ | ✅ Source→Deploy chain |
| **Trust Scoring** | ❌ | ❌ | ✅ Real-time scores |
| **SMS/Voice Alerts** | ❌ | ❌ | ✅ Get called on failures |
| **SMSLY Ecosystem** | ❌ | ❌ | ✅ Hosting + Identity |

## Architecture

```
smsly-code/
├── smsly-git/          # Rust — Git engine (libgit2)
├── smsly-code-api/     # Go — REST API server (Fiber)
├── smsly-code-web/     # Next.js 14 — Frontend
├── docker/             # Docker Compose + production setup
├── deploy/             # Kubernetes manifests
├── docs/               # Documentation (mdbook)
└── contracts/          # Integration contracts between components
```

## Tech Stack

- **Git Engine:** Rust + libgit2 (memory-safe, high-performance)
- **API Server:** Go + Fiber (fast HTTP, excellent concurrency)
- **Frontend:** Next.js 14 + TypeScript + Tailwind CSS
- **Database:** PostgreSQL 16
- **Search:** Meilisearch
- **Cache/Queue:** Redis 7
- **Object Storage:** MinIO / S3
- **AI:** Gemini API
- **Deployment:** Docker Compose + Kubernetes

## Quick Start

```bash
git clone https://github.com/SMSLYCLOUD/smsly-code.git
cd smsly-code
cp .env.example .env
make dev-up
# Visit http://localhost:3000
```

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions.

## License

MIT License — see [LICENSE](LICENSE) for details.

Built with ❤️ by SMSLY Cloud
