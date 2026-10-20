# Hookline

**A self-hostable webhook delivery platform, built one version at a time, to learn backend engineering from first principles.**

Companies like Stripe, GitHub and Shopify all need the same thing: when something happens in their system, they must tell their customers' servers about it reliably, securely, in order, with retries, and with a log of every attempt. Teams pay for products like Svix and Hookdeck to do this. Hookline is my own version of that product.

I picked it on purpose. A webhook platform can't be faked with CRUD:

- It **receives** traffic (an HTTP API that other backends call).
- It **sends** traffic to servers it doesn't control, which are slow, down, or hostile.
- It must not lose a message when a process crashes, deploys, or gets killed.
- It is multi-tenant, so one noisy customer must not slow down the rest.
- Every one of the 30 topics in the course is a real need of the product, not a bolt-on.

## How this project is built

This is a learning project with a strict method:

1. **The course is the foundation.** Every version maps to one or more videos from Sriniously's [Backend from First Principles](https://www.youtube.com/@sriniously) playlist, in playlist order. Watch the video first, then build.
2. **Versions, not features.** Hookline grows from `v0.0` (no code) to `v2.0` (a production-grade platform). Each version starts with a pain the previous version can't handle. That is how real software advances.
3. **Understand before code.** A mentor agent (instructions in [`CLAUDE.md`](./CLAUDE.md)) explains the why, the analogy and the design before any code is written. I write the code; the agent reviews it.
4. **Leave evidence.** Every version ships a learning-log entry, at least one Architecture Decision Record (ADR), a "break it" experiment, tests, and a git tag.

## Map

| Document | What it's for |
|---|---|
| [`ROADMAP.md`](./ROADMAP.md) | Every version, the videos it covers, what to build, and how to know it's done |
| [`CLAUDE.md`](./CLAUDE.md) | Instructions for the mentor agent that teaches me through the project |
| [`docs/architecture.md`](./docs/architecture.md) | Where the system is heading (the v2.0 picture) and the tech stack |
| [`docs/curriculum-map.md`](./docs/curriculum-map.md) | Video → version lookup table |
| [`docs/adr/`](./docs/adr/) | Architecture Decision Records |
| [`docs/learning-log/`](./docs/learning-log/) | One entry per version, in my own words |
| [`CHANGELOG.md`](./CHANGELOG.md) | What changed in each version |

## Current version

`v0.0` (Orientation). See [ROADMAP.md](./ROADMAP.md#v00--orientation).
