## Context

The EasyJet application currently has three areas requiring attention:

1. **Naming inconsistency**: The codebase uses "git logs" terminology for what are actually git commits. This creates confusion and doesn't align with standard git terminology.

2. **Unbounded data growth**: Project runs accumulate indefinitely with no cleanup mechanism, potentially causing database bloat and performance degradation over time.

3. **Missing static file serving**: The web server only serves API endpoints, requiring a separate server or manual deployment for the compiled frontend application.

Current architecture:
- Go backend with GORM for database operations
- SQLite database with tables: `runs`, `run_stages`, `run_git_logs`
- Vue.js frontend served separately
- Server runs HTTP API on configurable address

## Goals / Non-Goals

**Goals:**
- Rename all "git logs" references to "git commits" consistently across backend, frontend, and documentation
- Implement automatic run rotation with configurable retention policy
- Add static file serving with SPA routing support to the Go HTTP server
- Maintain backward compatibility during migration where possible
- Provide database migration for table rename

**Non-Goals:**
- Changing the git data collection logic (only renaming)
- Implementing complex retention policies (e.g., time-based, size-based)
- Adding CDN or advanced caching layers
- Modifying frontend build process

## Decisions

### Decision 1: Database migration approach for table rename

**Chosen approach:** Use GORM's AutoMigrate with explicit table rename in migration function.

**Rationale:** GORM doesn't support table renames directly in AutoMigrate. We'll create a manual migration step that executes raw SQL `ALTER TABLE run_git_logs RENAME TO run_git_commits` before running AutoMigrate.

**Alternatives considered:**
- Create new table and copy data: More complex, requires cleanup of old table
- Dual-write during transition: Overkill for internal tool
- Manual migration script: Less automated, error-prone

### Decision 2: Retention policy configuration

**Chosen approach:** Add `RetentionCount` field to Project entity with default value of 10.

**Rationale:** Allows per-project customization while providing sensible defaults. Stored in database for runtime access without config file reload.

**Alternatives considered:**
- Global config only: Less flexible for projects with different needs
- Environment variable: Harder to change per-project
- API-configurable: Adds complexity, retention is usually stable

### Decision 3: Rotation trigger timing

**Chosen approach:** Trigger rotation after successful run completion in the same worker goroutine.

**Rationale:** Simple, synchronous with run completion. Failed runs are preserved for debugging. No additional scheduler complexity needed.

**Alternatives considered:**
- Separate scheduled job: More complex, delayed cleanup
- Manual-only: Easy to forget, unbounded growth between manual runs
- Async trigger: Adds complexity for minimal benefit

### Decision 4: Static file serving implementation

**Chosen approach:** Use `http.FileServer` with custom handler for SPA fallback.

**Rationale:** Standard library solution, well-tested, minimal dependencies. Custom handler intercepts 404s for SPA routes while preserving API routes.

**Alternatives considered:**
- Third-party static server package: Adds dependency for simple need
- Embed frontend in binary: Larger binary, less flexible deployment
- Separate nginx: Operational complexity, defeats "single binary" goal

### Decision 5: Cache headers strategy

**Chosen approach:** Differentiate between index.html (no-cache) and asset files (long-term cache).

**Rationale:** Standard SPA deployment pattern. Ensures fresh HTML while allowing aggressive asset caching.

## Risks / Trade-offs

**[Risk] Database migration failure** → Mitigation: Wrap migration in transaction, provide rollback script. Test migration on copy of production data.

**[Risk] Breaking API change for git_commits** → Mitigation: Update frontend and backend together in same deployment. Document breaking change clearly.

**[Risk] Accidental data loss from rotation** → Mitigation: Log all deleted runs. Start with conservative default (10 runs). Add dry-run mode for testing.

**[Risk] SPA routing conflicts with API routes** → Mitigation: Explicitly check for `/api/` prefix before applying fallback. Test all API endpoints after implementation.

**[Trade-off] Per-project retention adds complexity** → Acceptable trade-off for flexibility. Default value handles simple cases.

**[Trade-off] Synchronous rotation slows run completion** → Acceptable for typical run counts (<100). Can be made async later if needed.

## Migration Plan

1. **Pre-deployment:**
   - Backup database
   - Review and test migration script on staging

2. **Deployment steps:**
   - Deploy backend with migration code
   - Run database migration (automatic on startup)
   - Deploy frontend with updated `git_commits` references
   - Verify API responses and UI rendering

3. **Rollback strategy:**
   - Restore database from backup
   - Deploy previous backend version
   - Redeploy previous frontend version

4. **Post-deployment:**
   - Monitor error logs for any `git_logs` references missed
   - Verify rotation is working after first successful runs
   - Confirm static files are served correctly

## Open Questions

1. Should we add an API endpoint to configure retention policy per project, or is config file sufficient?
2. Do we need to preserve a minimum number of failed runs regardless of rotation policy?
3. Should static files path support environment variable substitution for containerized deployments?
