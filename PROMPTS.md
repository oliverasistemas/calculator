# Prompts

The assignment asks to share the prompts used while building this project. All work was done
with [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Prompts are reproduced
verbatim; routine follow-ups ("run the tests", "fix that") are omitted.

## Initial build

<!-- TODO: paste the prompts used to scaffold and build the project. The session transcript
     was not available on the machine where the review/hardening pass below was done. -->

## Review and hardening pass

After the initial build, a review session was run against the assignment requirements.

**Prompt 1** — full review (the assignment text was pasted along with):

> Review project Objective [...assignment text...] this were the instructions

This produced a gap analysis against the deliverables list (missing coverage report, missing
prompts file), found an `int64` overflow bug in the backend's expression formatter, and flagged
smaller issues (untested `useCalculator` hook, CORS credentials with wildcard origins, a
decimal-point input edge case).

**Prompt 2** — acting on part of the review:

> Lets do coverage reports and prompts

This added the `make coverage*` targets, the coverage section in the README, and this file.

**Prompt 3** — closing the biggest coverage gap found in the review:

> add test for useCalculator

This added `app-web/src/hooks/useCalculator.test.ts` (20 tests mocking the API client),
bringing the hook to 100% coverage and the frontend total from 31% to 74%.
