# Supervision

How the lecturer (Antigravity) and the supervisor (Claude Code) work together on each version. You carry messages between them; the files here are the channel.

```
supervisor ── brief vX.Y ──► lecturer teaches, you build on version/vX.Y-name
    ▲                                              │
    │                               lecturer writes reports/vX.Y.md
    │                                              │
    └──────── review vX.Y ◄─────────────────────────┘
                 │
     APPROVED → merge, tag vX.Y.0, brief next version
     POLISH / REDO → lecturer runs a polish round, review again
```

## What you say to each agent

| Step | Where | You type |
|---|---|---|
| 1 | Claude Code | `brief v0.1` |
| 2 | Antigravity | `/start-version` |
| 3 | Antigravity | work through the version with the lecturer |
| 4 | Antigravity | `/submit-report` (then push the branch) |
| 5 | Claude Code | `review v0.1` |
| 6 | You | if approved: merge, `git tag v0.1.0`, push tags |

## Files

| Path | Written by | Contains |
|---|---|---|
| `STATUS.md` | both | Current version and state |
| `briefs/vX.Y.md` | supervisor | What to teach and how it will be judged |
| `reports/vX.Y.md` | lecturer | What was taught and how the learner did |
| `reviews/vX.Y.md` | supervisor | Verdict and what to fix or carry forward |
| `templates/` | — | Shapes for the three files above |

## States

`briefed` → `in-progress` → `reported` → `approved` | `polish` | `redo`
