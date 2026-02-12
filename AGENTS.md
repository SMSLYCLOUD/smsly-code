# Agent Instructions

This repository is designed to be built by a collaborative swarm of AI agents.

## How to Contribute

1.  **Pick a Plan**: Navigate to the `plans/` directory and select a plan file (e.g., `01_git_protocol.md`).
2.  **Read the Mission**: Understand the specific goals, tech stack, and verification steps for that agent.
3.  **Execute**: Implement the features described.
    *   Respect the existing monorepo structure.
    *   Update `contracts/` if you introduce new cross-service APIs.
    *   Add tests for your new features.
4.  **Verify**: Ensure the verification steps listed in the plan pass.
5.  **Update**: If you complete a plan, mark it as done or update the plan file with new learnings.

## Project Structure

*   `smsly-git/`: Rust-based Git engine.
*   `smsly-code-api/`: Go-based REST API.
*   `smsly-code-web/`: Next.js frontend.
*   `docker-compose.dev.yml`: Local development environment.

## Communication

If you need to change a contract between services (e.g., API response format), check `contracts/` first to ensure you don't break other agents' work.
