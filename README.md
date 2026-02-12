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

## Roadmap & Agent Plans

This project is built by a swarm of specialized agents. We have defined 10 distinct mission plans to reach feature parity with GitHub:

1.  **[Core Git Protocol](plans/01_git_protocol.md):** Smart HTTP, SSH, Hooks.
2.  **[Identity & Auth](plans/02_auth_identity.md):** OAuth, RBAC, Organizations, GPG Keys.
3.  **[Social & Discovery](plans/03_repo_social.md):** Forks, Stars, Wiki, Discussions.
4.  **[Code Review](plans/04_code_review.md):** Pull Requests, Diffs, Merge Strategies.
5.  **[Project Management](plans/05_project_management.md):** Issues, Kanban, Milestones.
6.  **[CI/CD Actions](plans/06_cicd_actions.md):** SMSLY Actions, Runners, Secrets.
7.  **[Security & MIP](plans/07_security_mip.md):** Integrity Stamps, Dependency Graph, Secret Scanning.
8.  **[Frontend Foundation](plans/08_frontend_foundation.md):** Design System, Command Palette, WebSockets.
9.  **[Frontend Features](plans/09_frontend_features.md):** Monaco Editor, Blame, Graphs.
10. **[Infrastructure](plans/10_infra_search.md):** Search, Scaling, K8s, Observability.

## License

MIT License — see [LICENSE](LICENSE) for details.

Built with ❤️ by SMSLY Cloud
