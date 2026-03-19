## Why

The codebase currently uses "git logs" terminology which is semantically incorrect — we're storing git commits, not logs. Additionally, the system needs run data rotation to prevent unbounded growth, and the web server must serve the compiled frontend static files for production deployment.

## What Changes

- **Rename "git logs" to "git commits"** throughout the codebase:
  - Entity `ProjectRunGitLogs` → `ProjectRunGitCommits`
  - Database table `run_git_logs` → `run_git_commits`
  - Field `GitLogs` → `GitCommits` in `ProjectRun` entity
  - All references in services, repositories, handlers, and frontend code
  - **BREAKING**: API response field changes from `git_logs` to `git_commits`
  - **BREAKING**: Database table name changes require migration

- **Add run data rotation capability**:
  - Automatic cleanup of old project runs based on configurable retention policy
  - Rotation triggered on schedule or after successful new run completion
  - Configurable retention count (default: keep last N runs per project)

- **Add static file serving to web server**:
  - Serve compiled frontend assets from configured directory
  - Enable SPA routing with fallback to index.html
  - Configurable static files path in server configuration

## Capabilities

### New Capabilities
- `run-rotation`: Automatic rotation and cleanup of old project runs based on retention policy
- `static-serving`: HTTP static file serving for compiled frontend assets with SPA routing support

### Modified Capabilities
- `git-commits`: Renaming from `git-logs` - all entity names, database tables, and API fields change from `git_logs` to `git_commits`

## Impact

- **Database**: Table `run_git_logs` renamed to `run_git_commits` requires migration script
- **API**: Breaking change — response field `git_logs` becomes `git_commits`
- **Frontend**: All references to `git_logs` must update to `git_commits`
- **Documentation**: All business docs, DFD diagrams, and entity descriptions need updates
- **Configuration**: New config field for static files directory and run retention policy
