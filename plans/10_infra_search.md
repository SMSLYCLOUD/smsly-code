# Agent 10: Infrastructure, Search & Analytics (Scale)

**Mission:** Build the robust backend infrastructure to index code for search, manage scaling, backups, and observability. This is the "Engine Room".

## Goals

1.  **Code & Issue Search (Meilisearch/Elastic)**
    *   **Goal:** Instant search across *all* repositories.
    *   Index code (symbol-aware or trigram-based).
    *   Index Issues, PRs, Wiki, Users, Orgs.
    *   Implement faceted search (language, repo, author, date).

2.  **Observability & Logging**
    *   Implement distributed tracing (OpenTelemetry) across Go, Rust, and Next.js.
    *   Centralize logs (Loki/Elastic).
    *   Create Grafana dashboards for API latency, Git performance, and Error rates.

3.  **Database Scaling & Caching**
    *   Implement read-replicas for PostgreSQL.
    *   Use Redis Cluster for caching frequently accessed data (repo metadata, permissions).
    *   Evaluate database partitioning/sharding strategies for huge repos.

4.  **Kubernetes Deployment (Helm)**
    *   Create Helm charts for `smsly-code` (api, git, web, actions-runner).
    *   Implement Horizontal Pod Autoscaling (HPA) based on CPU/Memory/Requests.
    *   Implement Ingress Controller (Nginx/Traefik) with TLS termination (Cert-Manager).

5.  **Backups & Disaster Recovery**
    *   Implement automated backups for Postgres (WAL-G/Barman).
    *   Implement incremental backups for Git repositories (restic/rclone to S3).
    *   Test restore procedures regularly.

## Tech Stack
*   **Search:** Meilisearch (v1.6+) or Elasticsearch
*   **Monitoring:** Prometheus, Grafana, OpenTelemetry, Jaeger
*   **Infrastructure:** Kubernetes, Helm, Terraform/Pulumi (AWS/GCP/Azure)
*   **Database:** PostgreSQL (Primary/Replica), Redis (Cluster)

## Verification
*   Searching for `fn main` returns relevant Rust files.
*   Grafana dashboard shows real-time request counts.
*   Deleting a pod causes Kubernetes to restart it automatically.
*   Backups are successfully uploaded to S3 daily.
