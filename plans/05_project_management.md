# Agent 5: Project Management & Issues (Go)

**Mission:** Build the project management suite (Issues, Kanban, Milestones) to empower teams to plan and track work effectively.

## Goals

1.  **Issues Tracking**
    *   [x] Implement standard CRUD operations for Issues (`POST /repos/{owner}/{repo}/issues`).
    *   [x] Support Markdown rendering (GFM) in issue descriptions.
    *   [ ] Implement `labels` (color-coded tags), `milestones`, and `assignees`.
    *   [ ] Implement cross-references (`#123` links to issue/PR).

2.  **Projects (Kanban/Tables)**
    *   [ ] Implement "Projects" at the User/Org level (not just Repo).
    *   [ ] Support custom fields (Status, Priority, Size).
    *   [ ] Views: Kanban Board (Drag & Drop), Table (Sort/Filter), Roadmap (Timeline).

3.  **Insights & Reporting**
    *   [ ] Implement Velocity charts, Burn-up/Burn-down charts.
    *   [ ] Track contribution graphs (commits/issues/PRs per user).

4.  **Mentions & Notifications**
    *   [ ] Implement `@mention` parsing in comments/descriptions.
    *   [ ] Send email notifications (SMTP/SendGrid).
    *   [ ] Deliver real-time web notifications (WebSockets).
    *   [ ] Allow users to manage notification preferences (`watch`, `participating`, `mention`).

## Tech Stack
*   **Language:** Go (Fiber)
*   **Database:** PostgreSQL (Issues, Comments, Projects)
*   **Frontend:** React DnD (for Kanban), Recharts (for Insights)

## Verification
*   User can create an issue, assign labels, and add it to a Milestone.
*   Moving a card on the Kanban board updates the issue status.
*   Mentioning `@user` triggers a notification for them.
*   Closing an issue via commit message (`Fixes #123`) works.
