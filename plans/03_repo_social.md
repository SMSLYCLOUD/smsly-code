# Agent 3: Repository Social & Discovery

**Mission:** Build the social graph and repository features that enable exploration, forking, and knowledge sharing.

## Goals

1.  **Forking & Network Graph**
    *   Implement `POST /repos/{owner}/{repo}/forks` to create a copy of a repository.
    *   Maintain a "Network" model to track the ancestry of forks (Parent/Child).
    *   Implement logic to prevent circular forks.
    *   Optimize fork creation (shallow copy or alternating objects if using libgit2 advanced features).

2.  **Social Interactions**
    *   **Stars**: `PUT /user/starred/{owner}/{repo}`.
    *   **Watch**: `PUT /repos/{owner}/{repo}/subscription` (Ignoring, Releases only, Participating).
    *   **Follow**: `PUT /user/following/{username}`.
    *   Update activity feed (Timeline) when these events occur.

3.  **Wiki & Pages**
    *   Implement a separate git repository for each wiki (`repo.wiki.git`).
    *   Serve markdown files as rendered HTML (`/repos/{owner}/{repo}/wiki`).
    *   (Optional) Support static site hosting (Pages) via MinIO/S3.

4.  **Discussions & Topics**
    *   Implement discussion forums (`/repos/{owner}/{repo}/discussions`).
    *   Support categories, locking, pinning, and marking answers.
    *   Implement repository Topics (tags) for discovery (`/topics/{topic}`).

5.  **Trending & Explore**
    *   Implement algorithms to rank repositories by stars/forks over time (Trending).
    *   Create curated collections.

## Tech Stack
*   **Language:** Go (Fiber)
*   **Database:** PostgreSQL (GORM) + Redis (ZSET for trending)
*   **Search:** Meilisearch (for topics/repos)

## Verification
*   User can fork a repository and see it listed under their account.
*   The original repository shows the fork count incremented.
*   Starring a repository adds it to the user's "Stars" tab.
*   Wiki pages can be created, edited, and viewed.
