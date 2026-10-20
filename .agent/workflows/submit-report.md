---
description: Write the version report and hand it to the supervisor for review
---

1. Confirm with the learner that every acceptance criterion in the brief is met. Run `gofmt -l .`, `go vet ./...` and `go test ./...` (skip if there's no Go code yet) and show the output.
2. Fill `supervision/templates/report.md` into `supervision/reports/vX.Y.md`. Record the learner's check-question answers verbatim. Keep it under ~80 lines.
   - If this is a polish round, append a "Polish round N" section instead of rewriting the report.
3. Set `supervision/STATUS.md` state to `reported`, with next action "In Claude Code, run `review vX.Y`".
4. Commit everything on the version branch and push it.
5. Tell the learner: "Ask the supervisor: review vX.Y".
