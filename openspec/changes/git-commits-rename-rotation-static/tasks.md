## 1. Database Migration Setup

- [x] 1.1 Create database migration function in `internal/adapter/repository/gorm/database.go`
- [x] 1.2 Add SQL statement to rename `run_git_logs` to `run_git_commits`
- [x] 1.3 Wrap migration in transaction with rollback on error
- [x] 1.4 Add migration logging for debugging

## 2. Backend Entity Renaming

- [x] 2.1 Rename `ProjectRunGitLogs` struct to `ProjectRunGitCommits` in `internal/core/entity/project_run.go`
- [x] 2.2 Update `TableName()` to return `run_git_commits`
- [x] 2.3 Rename `GitLogs` field to `GitCommits` in `ProjectRun` entity
- [x] 2.4 Update JSON tag from `git_logs` to `git_commits`
- [x] 2.5 Update GORM foreign key references

## 3. Repository Layer Updates

- [x] 3.1 Update `SetProjectRunGitLogs` method to `SetProjectRunGitCommits` in `internal/adapter/repository/gorm/project_run.go`
- [x] 3.2 Update all `Preload("GitLogs")` calls to `Preload("GitCommits")`
- [x] 3.3 Update all `GitLogs` field accesses to `GitCommits`
- [x] 3.4 Update sorting functions to use `GitCommits`
- [x] 3.5 Update interface method signature in `internal/core/port/interfaces.go`

## 4. Service Layer Updates

- [x] 4.1 Update `run_project.go` to use `SetProjectRunGitCommits`
- [x] 4.2 Update entity mapping from `ProjectRunGitLogs` to `ProjectRunGitCommits`
- [x] 4.3 Update error message from "save git logs" to "save git commits"

## 5. Frontend TypeScript Updates

- [x] 5.1 Rename interface `ProjectRunGitLog` to `ProjectRunGitCommit` in `frontend/src/pages/ProjectRun.vue`
- [x] 5.2 Update `git_logs` property to `git_commits` in `ProjectRun` interface
- [x] 5.3 Update template references from `run.git_logs` to `run.git_commits`
- [x] 5.4 Update `v-for` loop to use `git_commits`

## 6. Run Rotation Implementation

- [x] 6.1 Add `RetentionCount` field to Project entity in `internal/core/entity/project.go`
- [x] 6.2 Create rotation service function in `internal/core/service/run_rotation.go`
- [x] 6.3 Implement `RotateProjectRuns` method with retention logic
- [x] 6.4 Add rotation trigger after successful run completion in `run_project.go`
- [x] 6.5 Add manual rotation API endpoint `DELETE /api/projects/:id/runs/rotate`
- [x] 6.6 Add rotation handler in `internal/adapter/handler/httpapi/`
- [x] 6.7 Update database schema to include `retention_count` in projects table

## 7. Static File Serving Implementation

- [x] 7.1 Add `StaticFilesPath` field to server config in `config/config.go` (optional, no default value)
- [x] 7.2 Update `example.toml` with static files configuration example
- [x] 7.3 Create static file handler in `internal/adapter/handler/httpapi/static.go`
- [x] 7.4 Add conditional static handler registration - only if `StaticFilesPath` is configured
- [x] 7.5 Implement SPA fallback routing (serve index.html for unknown routes)
- [x] 7.6 Add cache headers for static assets
- [x] 7.7 Exclude `/api/*` routes from static file handling
- [x] 7.8 Register static handler in HTTP server setup in `cmd/server/main.go` with conditional check

## 8. Documentation Updates

- [x] 8.1 Update `docs/business/entity.md` to rename `ProjectRunGitLogs` to `ProjectRunGitCommits`
- [x] 8.2 Update `docs/business/screens.md` to use `git_commits` field name
- [x] 8.3 Update `docs/business/dfd.md` to rename data store references
- [x] 8.4 Update `docs/business/rules.md` to use new table name
- [x] 8.5 Update `docs/go/components.md` to use `ProjectRunGitCommits`
- [x] 8.6 Update `docs/product.md` to use "git commits" terminology
- [x] 8.7 Update ER diagrams in entity documentation

## 9. Testing and Verification

- [x] 9.1 Run `task go:test` to verify all tests pass
- [x] 9.2 Run `task go:lint` to verify no linting errors
- [x] 9.3 Run `task build:server` to verify build succeeds
- [ ] 9.4 Test database migration on test database
- [x] 9.5 Verify API response uses `git_commits` field
- [ ] 9.6 Test rotation with multiple runs exceeding retention
- [ ] 9.7 Test static file serving with built frontend
- [ ] 9.8 Verify SPA routing works for frontend routes

## 10. Final Cleanup

- [x] 10.1 Search codebase for any remaining `git_logs` references
- [x] 10.2 Search codebase for any remaining `GitLog` references
- [x] 10.3 Update any comments or documentation strings missed
- [x] 10.4 Run `task go:tidy` to clean up dependencies
