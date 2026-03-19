## RENAMED Requirements

### Requirement: Git commits entity naming

**FROM:** `ProjectRunGitLogs`

**TO:** `ProjectRunGitCommits`

**Reason:** The term "git logs" is semantically incorrect — the system stores individual git commits (hash, subject), not logs. This rename improves code clarity and aligns with git terminology.

**Migration:**

- Rename entity `ProjectRunGitLogs` to `ProjectRunGitCommits`
- Rename database table `run_git_logs` to `run_git_commits`
- Rename field `GitLogs` to `GitCommits` in `ProjectRun` entity
- Update API response field from `git_logs` to `git_commits`
- Update all frontend references from `git_logs` to `git_commits`
- Update all GORM `Preload("GitLogs")` calls to `Preload("GitCommits")`
- Update all struct field accesses `.GitLogs` to `.GitCommits`
- Update all type references `entity.ProjectRunGitLogs` to `entity.ProjectRunGitCommits`
- Update all slice types `[]entity.ProjectRunGitLogs` to `[]entity.ProjectRunGitCommits`
- Update all function parameters and return types using the old entity name

### Requirement: Git commits API response field

**FROM:** `git_logs`

**TO:** `git_commits`

**Reason:** Consistency with entity rename and accurate representation of data (commits, not logs).

**Migration:** Frontend must update all references to use `git_commits` field in API responses.

### Requirement: GORM preload strings

**FROM:** `Preload("GitLogs")`

**TO:** `Preload("GitCommits")`

**Reason:** GORM uses string-based preload references that must match the struct field name.

**Migration:** Search and replace all `Preload("GitLogs")` and `Preload("GitLogs.")` occurrences in repository layer.
