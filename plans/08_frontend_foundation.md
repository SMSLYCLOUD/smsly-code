# Agent 8: Frontend Foundation & Design System (Next.js)

**Mission:** Establish the visual language, layout, and core interactivity of the SMSLY Code web platform. This is the "face" of the product.

## Goals

1.  **Design System (Components)**
    *   Build a reusable component library (`smsly-ui`) using Tailwind CSS + Radix UI.
    *   Components: Button, Input, Modal, Dropdown, Tabs, Badge, Avatar, Select, Tooltip.
    *   Dark Mode / Light Mode support (via `next-themes`).
    *   Accessibility (ARIA) compliance.

2.  **Global Layout & Navigation**
    *   Implement responsive Header (Search, Notifications, Profile, Create Repo).
    *   Implement Sidebar (Context-aware: Global vs Repo vs Org vs Settings).
    *   Implement Breadcrumbs (`owner / repo / tree / branch / path`).

3.  **Real-Time Infrastructure**
    *   Implement WebSocket connection (Socket.IO / Phoenix Channels via Go API).
    *   Use WebSockets for:
        *   Notification counts (bell icon).
        *   CI/CD log streaming.
        *   PR comment updates (avoid page refresh).

4.  **Command Palette (`Cmd+K`)**
    *   Implement a global command menu (cmdk).
    *   Search repositories, jump to files, execute commands (create issue, etc.).

5.  **Keyboard Shortcuts**
    *   Implement `?` modal to show shortcuts.
    *   `g c` (Go to Code), `g i` (Go to Issues), `t` (File finder), `s` (Search).

## Tech Stack
*   **Framework:** Next.js 14 (App Router)
*   **Styling:** Tailwind CSS, Radix UI Primitives, Lucide Icons
*   **State:** React Query (TanStack Query), Zustand (Global store)
*   **WebSockets:** `socket.io-client`

## Verification
*   Dark mode toggle persists preference.
*   `Cmd+K` opens command palette and can navigate to a repo.
*   Sidebar collapses on mobile.
*   Notifications appear without page refresh.
