## 1. Database and Entity Layer

- [x] 1.1 Add `CronEnabled bool` field to `Project` entity in `internal/core/entity/project.go`
- [x] 1.2 Add `CronSchedule string` field to `Project` entity (plain string, not pointer)
- [x] 1.3 Set default values: `CronEnabled=false`, `CronSchedule=""` in GORM model
- [x] 1.4 Update GORM repository to handle new fields (auto-migration will handle DB schema)
- [ ] 1.5 Create database migration test to verify columns are added correctly
- [ ] 1.6 Test migration on both SQLite and PostgreSQL

## 2. Backend API Implementation

- [x] 2.1 Update API request/response models to include `cron_enabled` and `cron_schedule` fields
- [x] 2.2 Add `go-co-op/gocron` dependency to `go.mod`
- [x] 2.3 Add cron expression validation logic (use `gocron.CronExpressionParser`)
- [x] 2.4 Update `create_project` handler to accept and validate cron fields
- [x] 2.5 Update `update_project` handler to accept and validate cron fields
- [x] 2.6 Update `project` handler to return cron fields in response
- [x] 2.7 Add validation error responses for invalid cron expressions
- [ ] 2.8 Write unit tests for cron validation logic
- [ ] 2.9 Write integration tests for API endpoints with cron fields

## 3. Event Queue Adapter

- [x] 3.1 Create event queue adapter in `internal/adapter/scheduler/events.go`
- [x] 3.2 Define `SchedulerEvent` struct (Type, ProjectID, Schedule, Enabled)
- [x] 3.3 Define `EventType` enum (Created, Updated, Deleted)
- [x] 3.4 Implement channel-based event queue (`chan SchedulerEvent`)
- [x] 3.5 Implement event publisher interface for service layer
- [ ] 3.6 Write unit tests for event queue adapter

## 4. Scheduler Service Implementation

- [x] 4.1 Create scheduler service in `internal/adapter/handler/scheduler/scheduler.go`
- [x] 4.2 Initialize `gocron` scheduler with proper configuration
- [x] 4.3 Implement scheduler startup logic (start with application)
- [x] 4.4 Implement event queue consumer (read from channel)
- [x] 4.5 Implement job registration using `gocron.Cron()` method
- [x] 4.6 Implement job update logic (remove old, add new)
- [x] 4.7 Implement job removal when project is deleted/disabled
- [x] 4.8 Implement job execution callback (create pending run)
- [x] 4.9 Integrate scheduler with event queue adapter
- [x] 4.10 Integrate scheduler with existing service layer
- [x] 4.11 Add scheduler logging (INFO for scheduled runs, ERROR for failures, DEBUG for events)
- [x] 4.12 Add graceful shutdown handling for scheduler
- [x] 4.13 Update `cmd/server/main.go` to start scheduler alongside worker
- [ ] 4.14 Write unit tests for scheduler logic
- [ ] 4.15 Write integration tests for end-to-end scheduled run creation

## 5. Frontend UI Implementation

- [x] 5.1 Add `cron_enabled` and `cron_schedule` fields to TypeScript interfaces
- [x] 5.2 Add cron schedule toggle switch to ProjectEditor form in `ProjectEditor.vue`
- [x] 5.3 Add cron schedule input field with proper validation
- [x] 5.4 Add cron format hint and examples (e.g., "0 5 \* \* \* = daily at 5:00 AM")
- [x] 5.5 Add frontend validation for cron expressions
- [x] 5.6 Implement toggle behavior (enable/disable without losing schedule)
- [x] 5.7 Update Project details page to display cron schedule status
- [ ] 5.8 Test form submission with cron fields
- [ ] 5.9 Test form validation for invalid cron expressions
- [ ] 5.10 Test toggle on/off behavior
- [ ] 5.11 Test responsive design on different screen sizes

## 6. Documentation Updates

- [x] 6.1 Update `docs/go/libs.md` to include `go-co-op/gocron`
- [x] 6.2 Update `docs/business/entity.md` to include `CronEnabled` and `CronSchedule` fields in Project entity
- [x] 6.3 Update `docs/business/screens.md` to include cron controls in Project Editor screen
- [x] 6.4 Update `docs/go/arch.md` to include scheduler and event queue components
- [ ] 6.5 Update `docs/go/components.md` to include scheduler and event queue adapter
- [ ] 6.6 Update API documentation with cron fields
- [ ] 6.7 Add user guide section on configuring cron schedules
- [ ] 6.8 Document limitation: server downtime misses scheduled runs
- [ ] 6.9 Document timezone behavior: uses server local time

## 7. Testing and Verification

- [x] 7.1 Run full test suite: `task go:test`
- [x] 7.2 Run linter: `task go:lint`
- [x] 7.3 Run frontend linter: `task ts:lint`
- [ ] 7.4 Manual testing: Create project with cron scheduling enabled
- [ ] 7.5 Manual testing: Verify scheduled run creation
- [ ] 7.6 Manual testing: Test cron schedule toggle on/off
- [ ] 7.7 Manual testing: Test cron schedule updates (event queue processing)
- [ ] 7.8 Manual testing: Test project deletion (job removal)
- [ ] 7.9 Manual testing: Test removing cron schedule
- [ ] 7.10 Performance test: Verify low resource usage on scheduler

## 8. Build and Deployment

- [x] 8.1 Build server: `task build:server`
- [x] 8.2 Build frontend: `task build:front`
- [ ] 8.3 Test full application: `task run:server` + `task run:front`
- [ ] 8.4 Verify database migration runs successfully on fresh install
- [ ] 8.5 Verify database migration runs successfully on existing installation
- [ ] 8.6 Create deployment notes for cron scheduling feature
