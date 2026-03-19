## ADDED Requirements

### Requirement: Run rotation configuration

The system SHALL provide configuration for run retention policy per project.

#### Scenario: Default retention policy
- **WHEN** no retention policy is configured for a project
- **THEN** system uses default retention count of 10 runs per project

#### Scenario: Custom retention policy
- **WHEN** a project has a custom retention policy configured
- **THEN** system keeps only the specified number of most recent runs

### Requirement: Automatic run rotation trigger

The system SHALL automatically rotate old runs when a new run completes successfully.

#### Scenario: Rotation after successful run
- **WHEN** a new project run completes with success status
- **THEN** system checks if run count exceeds retention policy
- **THEN** system deletes oldest runs exceeding the retention limit

#### Scenario: No rotation on failed run
- **WHEN** a project run fails
- **THEN** system does not trigger rotation (preserve failed run for debugging)

### Requirement: Run rotation cleanup

The system SHALL delete old runs and all associated data during rotation.

#### Scenario: Delete old run with all relations
- **WHEN** a run is selected for rotation
- **THEN** system deletes the run and all associated run_stages records
- **THEN** system deletes all associated run_git_commits records

### Requirement: Manual rotation API endpoint

The system SHALL provide an API endpoint to manually trigger run rotation for a project.

#### Scenario: Manual rotation request
- **WHEN** user sends DELETE request to `/api/projects/:id/runs/rotate`
- **THEN** system rotates runs exceeding retention policy
- **THEN** system returns count of deleted runs
