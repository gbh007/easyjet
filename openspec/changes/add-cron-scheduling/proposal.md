## Why

EasyJet currently supports manual project runs and background processing of pending runs, but lacks automated scheduling. Users must manually trigger each pipeline execution, which defeats the purpose of CI/CD automation. Adding cron-based scheduling will enable automatic periodic execution of projects, aligning with the MVP goal of supporting cron-based job scheduling.

## What Changes

- **New Field**: Add `CronSchedule` field to `Project` entity (nullable string)
- **Database**: New migration for cron schedule column in projects table
- **API**: Add cron schedule to create/update project endpoints
- **Frontend**: Add cron schedule input to project editor form
- **Scheduler**: New background service to check and trigger scheduled runs
- **Worker Integration**: Scheduler integrates with existing worker infrastructure
- **Documentation**: Update business docs, API docs, and user guides

## Capabilities

### New Capabilities

- `cron-scheduling`: Core cron scheduling functionality including database schema, scheduler service, and integration with existing run infrastructure
- `cron-ui`: Frontend components and API integration for configuring cron schedules in the project editor

### Modified Capabilities

- `project-management`: Extending project entity and APIs to include cron schedule configuration

## Impact

**Affected Components:**

- **Entity Layer**: `Project` struct gains `CronSchedule` field
- **Database Layer**: GORM repository needs migration and query updates
- **Service Layer**: New scheduler service, integration with existing `RunProject` logic
- **HTTP API**: Request/response models updated, validation rules
- **Frontend**: Project editor form, potential new display on project details
- **Worker**: Existing worker may need to consider cron-triggered runs differently

**Dependencies:**

- Requires robust cron expression parsing (existing Go `robfig/cron` or similar)
- Must integrate with existing `PendingProjectRuns` and `HandleRun` infrastructure
- Should respect existing logging and monitoring patterns

**Non-Goals (for this change):**

- Webhook-based triggers (separate change)
- Advanced scheduling features (timezones, complex schedules beyond cron)
- Per-stage scheduling (schedule applies to entire project pipeline)
