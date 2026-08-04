# Prompts

The assignment asks to share the prompts used while building this project. All work was done
with [Claude Code](https://docs.anthropic.com/en/docs/claude-code). The prompts below are
reproduced verbatim where possible, with minor contextual descriptions added to explain the
purpose of each interaction. Routine follow-ups such as "run the tests" or "fix that" are
omitted.

## Initial build

**Prompt 1** — scaffolding the initial full-stack application. The full assignment text was
pasted along with a reference to an existing Go project for structural guidance:

> Using similar structure to ~/inmatics/gcbapos/ — Build a full-stack calculator application
> with a React frontend and a backend microservice. The frontend should consume the backend API
> to perform basic and advanced arithmetic operations. Focus on clean design, maintainable code,
> and testable architecture. Follow the assignment requirements closely, including the expected
> API behavior, project structure, testing requirements, Docker setup, and documentation.
> [full assignment requirements pasted]

**Prompt 2** — correcting the Go module layout after the initial scaffold placed `go.mod` at
the repository root instead of inside the backend service:

> Move go.mod into services/api/ and update the backend project structure and any references
> that depend on the module location.

**Prompt 3** — fixing the Docker build after moving `go.mod`. The full `docker compose` error
output was pasted to provide the build context and failure details:

> We have a similar issue with go.mod location when using Docker Compose. The Dockerfile is
> currently expecting go.mod at the repository root, but the Go module now lives under
> services/api/. Can you update the Docker configuration and build context so the backend
> continues to build correctly?

**Prompt 4** — generating the calculator frontend using the `/frontend-design` slash command.
The goal was to create a polished calculator interface using Ant Design components while
keeping the UI consistent with the application's requirements.

**Prompt 5** — reviewing the generated frontend for potential improvements:

> Review the generated frontend and the overall calculator user experience. Do you see any
> improvements that could be made to the UI, component structure, accessibility, or code
> organization? Please make any improvements that are clearly beneficial without unnecessarily
> complicating the implementation.

**Prompt 6** — manually testing an important arithmetic edge case and asking Claude to verify
whether the application handled it correctly:

> How much is 9 divided by 0?
>
> Does the calculator currently handle division by zero correctly? Please review both the
> frontend behavior and backend API behavior, and make any changes needed to handle this case
> gracefully and consistently.

**Prompt 7** — adding structured application logging with support for different log levels:

> Can you create a logger for the backend that supports at least debug and info log levels?
> Keep the implementation simple and idiomatic for Go, and make sure the logger can be
> configured and used consistently throughout the application.

**Prompt 8** — refactoring the logging configuration into a dedicated types package:

> Extract the log level configuration into a dedicated config type that lives in a new types
> package. Update the existing logger and configuration code to use this type cleanly and
> consistently.

**Prompt 9** — adding environment variable support using `godotenv`:

> Add support for loading environment variables from a `.env` file using
> "github.com/joho/godotenv". Integrate it into the application configuration without
> breaking the existing environment variable behavior.

**Prompt 10** — reviewing the latest changes and addressing a redundant implementation detail:

> Review the recent changes and identify any issues or unnecessary code.
>
> Fix these issues, including the following:
>
> cfg.LogLevel.Level() is redundant

**Prompt 11** — checking whether an unused backend endpoint was still required:

> Is the backend GET operation used anywhere by the frontend or required by the assignment?
> Please verify whether it is actually needed. If it is unused and unnecessary, remove the
> endpoint and any related code and tests.

**Prompt 12** — repairing the backend test suite after the API refactoring:

> Review the backend test suite after the recent API changes. Fix any failing, outdated, or
> missing tests and make sure the backend tests accurately reflect the current implementation
> and expected behavior.

**Prompt 13** — generating the README based on the assignment requirements and documenting
the architectural and implementation decisions:

> This application is a take home assignment as part of a take home test for a job application.
> Can you help write the README.md with the choices made, architectural decisions, setup
> instructions, how to run the application, how to run the tests, and any other information
> that would be useful for someone reviewing the project?
>
> These were the instructions I received:
> [assignment text pasted]

**Prompt 14** — standardizing the Go version across the project configuration:

> Use Go 1.26 consistently throughout the project. Update the Go version anywhere it is
> specified, including go.mod, Dockerfiles, CI configuration, Makefile, and documentation,
> while ensuring the project remains internally consistent.

**Prompt 15** — adding development and quality-of-life build tooling:

> Can you create a Makefile with useful commands for developing, testing, building, and
> maintaining the project?
>
> Add support for running go vet and golangci-lint as part of the development workflow.
>
> Is there a single Makefile target that can run the backend tests, frontend tests, go vet,
> and linting together so I can use it as a general quality check before submitting the
> project?

**Prompt 16** — documenting the newly added Makefile commands:

> Add a section to README.md documenting the Makefile and the available commands. Include
> instructions for the most common development, testing, linting, and quality-check commands
> so that another developer can quickly understand how to use the project tooling.

## Review and hardening pass

After the initial implementation was complete, a separate review session was run against the
original assignment requirements to identify missing deliverables, potential bugs, and areas
where the implementation could be strengthened before submission.

**Prompt 1** — performing a comprehensive review against the original assignment. The full
assignment text was pasted along with the request:

> Review the project against the following objective and assignment requirements. Check the
> implementation, architecture, tests, documentation, and project deliverables carefully.
> Identify anything that is missing, incomplete, incorrect, or potentially problematic.
> Pay particular attention to requirements that may have been overlooked during implementation
> and to edge cases that could cause incorrect behavior.
>
> Objective: [...]
>
> These were the instructions:
> [assignment text pasted]

This produced a gap analysis against the deliverables list. The review identified a missing
coverage report and missing prompts file, found an `int64` overflow bug in the backend's
expression formatter, and flagged smaller issues including the untested `useCalculator` hook,
CORS credentials being enabled with wildcard origins, and a decimal-point input edge case.

**Prompt 2** — implementing the missing coverage and prompt documentation identified during
the review:

> Let's add the coverage reports and the prompts documentation. Please implement the required
> Makefile coverage targets, update the README with instructions for generating and viewing
> coverage reports, and create the prompts documentation required by the assignment.

This added the `make coverage*` targets, the coverage section in the README, and this file.

**Prompt 3** — closing the biggest frontend coverage gap identified during the review:

> Add comprehensive tests for the useCalculator hook. Cover the main success paths as well as
> error handling and relevant edge cases. Mock the API client as necessary and follow the
> existing frontend testing conventions. Make sure the tests provide meaningful coverage of
> the hook's behavior rather than simply increasing the coverage percentage.

This added `app-web/src/hooks/useCalculator.test.ts` with 20 tests mocking the API client,
bringing the hook to 100% coverage and increasing the overall frontend coverage from 31% to
74%.
