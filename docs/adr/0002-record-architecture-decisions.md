# ADR-0002: Record Architecture Decisions

- **Status:** accepted
- **Version:** v0.0
- **Date:** 2026-10-22

## Context

Hookline is being engineered from first principles across 26 distinct versions, progressing from raw sockets to a production-grade webhook delivery platform with persistence, queues, caching, rate limiting, and WebSockets.

Throughout this journey, critical technical decisions will be made—such as database technology, queue implementations (hand-built Postgres vs. specialized message brokers), API contracts, and concurrency models.

Without a structured method for documenting decisions:
- The context, constraints, and trade-offs that led to a decision are quickly forgotten.
- Developers may repeatedly debate resolved topics or mistakenly undo intentional architectural constraints.
- Code comments and Git commit messages provide fragmented, localized reasoning rather than cohesive architectural records.

## Options considered

1. **Rely on Git Commit Messages and Pull Request Discussions:**
   - *Pros:* Requires no extra files or templates.
   - *Cons:* Difficult to search and browse as a cohesive architectural history; loses visibility once PRs are merged.
2. **External Wiki / Confluence:**
   - *Pros:* Rich WYSIWYG editing and easy sharing with non-technical stakeholders.
   - *Cons:* Disconnected from the codebase; easily drifts out of sync with actual code changes; cannot be reviewed as part of git pull requests.
3. **Architecture Decision Records (ADRs) in the Git Repository (Chosen):**
   - *Pros:* Versioned right alongside the source code; reviewed in pull requests; text-based Markdown format; standardized structure (Context, Options, Decision, Consequences).

## Decision

We adopt Architecture Decision Records (ADRs) using a numbered markdown format (`docs/adr/NNNN-title.md`) based on the template in `docs/adr/0000-template.md`.

An ADR must be created whenever a decision:
- Introduces or removes a significant dependency or technology stack component.
- Changes the public API surface or versioning scheme.
- Alters data storage, concurrency, or communication patterns between subsystems.

## Consequences

- **What gets easier:**
  - Clear traceability of architectural rationale and trade-offs for future maintenance and senior engineering reviews.
  - Onboarding and code reviews have an authoritative source for design decisions.
- **What gets harder:**
  - Requires maintaining documentation discipline before merging architectural changes.
