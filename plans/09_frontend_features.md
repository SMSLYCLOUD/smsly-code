# Agent 9: Frontend Feature Implementation (Next.js)

**Mission:** Develop the sophisticated frontend experiences (Code Viewer, Blame, Diff Viewer) that make working with code a delight.

## Goals

1.  **Advanced Code Viewer (Monaco)**
    *   **Goal:** Read-only code editor experience.
    *   Integrate `monaco-editor` or `codemirror`.
    *   Syntax highlighting for 100+ languages (Tree-sitter preferred).
    *   "Go to definition" / "Find references" (Language Servers via WebSockets - advanced).

2.  **Blame & History Views**
    *   **Blame:** Implement line-by-line blame view (`git blame`).
    *   Hover on a line to see commit message, author, and timestamp.
    *   **History:** List commits for a specific file/folder.

3.  **File Finder (`t`) & Symbol Search**
    *   Implement "Fuzzy File Finder" (Cmd+P like behavior) for the current repository tree.
    *   Index symbols (functions/classes) for quick jump (requires backend indexing).

4.  **Complex Diff Viewer (PRs)**
    *   Side-by-side vs Unified view.
    *   Code folding (Collapse unchanged sections).
    *   "View File" context from diff.
    *   Image diffs (Swipe/Onion skin) for binaries.

5.  **Interactive Graphs (Charts)**
    *   **Goal:** Visualize repo activity.
    *   Contributors graph (Commits/LOC over time).
    *   Network graph (Forks visualization - D3.js/Cytoscape).
    *   Dependency graph visualization.

## Tech Stack
*   **Libraries:** `monaco-editor`, `cmdk`, `framer-motion`, `recharts`, `cytoscape`
*   **Styling:** Tailwind CSS

## Verification
*   Viewing a `.rs` file shows syntax highlighting.
*   Selecting lines creates a permalink (`#L10-L20`).
*   Blame view correctly attributes lines to commits.
*   Clicking `t` opens fuzzy file finder and jumps to file.
