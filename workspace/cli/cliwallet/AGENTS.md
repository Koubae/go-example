AGENTS
======


This file is the main entry point for AI coding agents working in this repository.

The canonical, agent-agnostic rules and documentation are defined under:

[.agents/](./.agents/)

Agents should prefer the instructions in [.agents/](./.agents/) over IDE- or vendor-specific configuration.

IDE-specific files such as .cursor/, .github/, .claude/, or similar should act only as adapters and point back to the shared rules in .agents/.

Do not duplicate rules across different agent or IDE configurations.
