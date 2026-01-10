# Fragtape

Fragtape is a self-hosted tool for **automatically generating CS2 highlight clips** and collecting fun group statistics.
It does not aim to be a performance analysis tool.

> [!IMPORTANT]
> **Work in Progress**
>
> Fragtape is under active development and **does not work yet**.
> Expect breaking changes and missing features.

Planned features include:

- Automatic detection of new CS2 matches
- Faceit integration
- Demo parsing to extract highlight moments
- Configurable highlight rules (kills, clutches, intervals, ...)
- Rendering clips using CS2 demos
- Group based settings and actions (e.g. Discord posting)
- Simple group stats (knife kills, team kills, ...)

## Development

### Quick Start

1. Install the tools listed in the [asdf file](./.tool-versions)
2. Install _make_.
3. Run `make setup` to install:

- Golang tools: _Air_, _Goose_, _Sqlc_, _Deadcode_
- Frontend dependencies

1. Install the git hook for code quality: `git config --local core.hooksPath .githooks/`
2. Copy `.env.example` -> `.env` and populate
3. Run database migrations: `make migrate`
4. Start the project `make watch`.

Endpoints:

- **Backend:** <http://localhost:3001>
- **Frontend:** <http://localhost:3000>

### High Level Architecture

- **Server**
  - User facing API and UI
  - Stores configuration and metadata
  - Serves highlight clips
- **Worker**
  - Detects new matches
  - Downloads demos
  - Parses demos and plans highlight segments
  - Adds segments to job queue
  - Processes finished clips (e.g. send discord messages)
- **Renderer**
  - Consumes jobs
  - Runs CS2 + tooling to generate clips
  - Uploads results to object storage
  - Horizontally scalable
