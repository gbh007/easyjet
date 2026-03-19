## Context

EasyJet's MVP requires cron-based job scheduling, but the current system only supports manual triggering of project runs. The architecture uses a worker service that polls for pending runs every second. We need to add a scheduler component that evaluates cron expressions and automatically creates pending runs when scheduled.

**Current State:**

- Projects have no scheduling configuration
- Worker polls every second for `pending=true` runs
- All runs are triggered manually via API
- No background scheduling infrastructure exists

**Constraints:**

- Must work on low-power hardware (Raspberry Pi)
- Single-node architecture (no distributed scheduling)
- Must integrate with existing `RunProject` and `HandleRun` flow
- SQLite and PostgreSQL support required

**Stakeholders:**

- Users: Need reliable automated scheduling
- DevOps: Simple deployment, no external dependencies (like Redis)

## Goals / Non-Goals

**Goals:**

- Add cron schedule configuration to projects
- Implement scheduler service that evaluates cron expressions
- Automatically create pending runs when schedule triggers
- Integrate seamlessly with existing worker infrastructure
- Provide UI for configuring cron schedules
- Support standard cron expression format (5 fields)

**Non-Goals:**

- Timezone support (use server local time)
- Complex scheduling (e.g., "every 5th Monday")
- Per-stage scheduling (schedule applies to entire pipeline)
- Webhook-based triggers (separate feature)
- Distributed/multi-node scheduling
- Backfill missed runs (if server was down)

## Decisions

### Decision 1: Use `go-co-op/gocron` Library

**Choice:** Use `github.com/go-co-op/gocron` for cron expression parsing and scheduling

**Rationale:**

- Feature-rich scheduler with clean API
- Supports standard cron format and human-readable scheduling
- Well-maintained with active development
- Provides built-in job management (start/stop/resume)
- Better error handling and logging support

**Alternatives Considered:**

- `robfig/cron/v3`: Industry standard but less feature-rich
- Custom implementation: Error-prone, reinventing the wheel
- `adhocore/gronx`: Good but less battle-tested

### Decision 2: Scheduler as Separate Background Service

**Choice:** Run scheduler as independent goroutine alongside existing worker

**Architecture:**

```
┌─────────────────┐
│   Scheduler     │──┐
│  (cron checker) │  │
└─────────────────┘  │
                     ▼
              ┌─────────────┐
              │  Database   │
              │ (pending    │
              │   runs)     │
              └─────────────┘
                     ▲
                     │
┌─────────────────┐  │
│     Worker      │──┘
│ (run executor)  │
└─────────────────┘
```

**Rationale:**

- Clean separation of concerns
- Scheduler creates pending runs, worker executes them (existing pattern)
- No changes needed to worker execution logic
- Easy to test independently

**Alternatives Considered:**

- Integrate cron check into existing worker loop: Would complicate worker, mix concerns
- Single unified service: Harder to reason about, test

### Decision 3: Two-Field Cron Configuration (Enable/Disable + Schedule)

**Choice:** Add two fields to `Project` entity:

- `CronEnabled bool` - toggle to enable/disable scheduling
- `CronSchedule string` - the cron expression (empty string means no schedule)

**Schema:**

```go
type Project struct {
    // ... existing fields ...
    CronEnabled    bool   `json:"cron_enabled" gorm:"column:cron_enabled;not null;default:false"`
    CronSchedule   string `json:"cron_schedule" gorm:"column:cron_schedule;type:text;not null;default:''"`
}
```

**Rationale:**

- **Explicit control**: Users can disable scheduling without losing the cron expression
- **Safe editing**: Can temporarily disable, then re-enable with same schedule
- **Clear intent**: `CronEnabled=false` clearly shows scheduling is disabled
- **Flexible**: Can have `CronSchedule` set but disabled (for testing, debugging)
- **Default safe**: New projects start with `CronEnabled=false` (no accidental runs)
- **UI toggle**: Simple switch to enable/disable without clearing schedule
- **No pointers**: Simpler code, no nil checks, empty string is clear "no value" indicator

**Behavior:**

- Scheduler only triggers runs when BOTH `CronEnabled=true` AND `CronSchedule!=""`
- UI shows toggle switch + cron expression input
- Disabling doesn't clear the schedule, just prevents execution
- Empty string `""` means "no schedule" (not NULL)

**Alternatives Considered:**

- Pointer to string (`*string`): Requires nil checks, adds complexity without benefit
- Nullable field: Unclear if NULL means "disabled" or "not set"
- Special string values (e.g., "disabled"): Hacky, error-prone
- Separate schedule table: Overkill for this use case

### Decision 4: Use gocron.Cron() with Dynamic Job Management

**Choice:** Use `gocron.Cron()` method with dynamic job registration/removal via event queue

**Implementation:**

```go
type Scheduler struct {
    scheduler  *gocron.Scheduler
    eventQueue chan SchedulerEvent
    jobs       map[uint]gocron.Job // project ID -> job mapping
}

type SchedulerEvent struct {
    Type      EventType // Created, Updated, Deleted
    ProjectID uint
    Schedule  string
    Enabled   bool
}

// Start event processor
go func() {
    for event := range s.eventQueue {
        s.handleEvent(event)
    }
}()
```

**Rationale:**

- **Dynamic updates**: Jobs are registered/updated/removed when configuration changes
- **No polling**: Scheduler doesn't need to check database every minute
- **Efficient**: Only active jobs are in memory, no wasted cycles
- **Clean architecture**: Event queue decouples API from scheduler
- **Channel-based**: Go's CSP model for concurrent communication

**Alternatives Considered:**

- Poll database every minute: Wasteful, doesn't scale, latency
- Use gocron singleton with tags: Less control, harder to manage
- Check all projects on each tick: Inefficient for large project counts

### Decision 5: No Concurrent Run Protection Needed

**Choice:** Skip checking for active runs when triggering scheduled runs

**Rationale:**

- **Single instance**: EasyJet runs as a single instance on one server
- **Single worker**: One worker processes runs sequentially
- **No race conditions**: No parallel execution possible
- **Simpler code**: No need for database queries to check active runs
- **Lower latency**: Direct run creation without additional checks

**Behavior:**

- Scheduler creates pending runs unconditionally when cron triggers
- Worker processes runs one at a time in FIFO order
- If a run is already pending/processing, new runs queue behind it
- Natural backpressure from single worker

**Alternatives Considered:**

- Check for active runs: Unnecessary complexity for single-instance system
- Limit queue size: Could lose scheduled runs, adds complexity
- Cancel existing runs: Unexpected behavior, could lose work

### Decision 6: API Backward Compatibility

**Choice:** Make cron fields optional in API requests with sensible defaults

**Request Model:**

```go
type Project struct {
    // ... existing fields ...
    CronEnabled    bool   `json:"cron_enabled"`
    CronSchedule   string `json:"cron_schedule"`
}
```

**Rationale:**

- **Non-breaking change**: Existing clients continue working
- **Defaults**: `CronEnabled=false`, `CronSchedule=""` for old clients
- **Frontend can add feature progressively**: Toggle and input added in UI update

**Behavior:**

- Old API clients: Fields default to disabled/empty (no behavior change)
- New API clients: Can set both fields explicitly
- Validation: Only validate `CronSchedule` if non-empty

## Risks / Trade-offs

### Risk 1: Server Downtime Misses Scheduled Runs

**Risk:** If EasyJet is stopped during a scheduled time, the run is missed

**Impact:** Users expect runs to happen even if server was temporarily down

**Mitigation:**

- Document this limitation clearly
- Future enhancement: Add "backfill" configuration option
- Acceptable for MVP (standard cron behavior)

### Risk 2: Timezone Confusion

**Risk:** Server local time may not match user expectations

**Impact:** Runs may trigger at unexpected times for users in different timezones

**Mitigation:**

- Document that cron uses server local time
- Future enhancement: Add timezone configuration
- Consider UTC as default in future versions

### Risk 3: Database Migration on Existing Installations

**Risk:** Adding column to existing databases may fail

**Impact:** Upgrade failures for users with existing data

**Mitigation:**

- GORM auto-migration handles this gracefully
- Test migration on SQLite and PostgreSQL
- Provide rollback script (add migration down function)

### Risk 4: Scheduler Performance with Many Projects

**Risk:** Checking hundreds of projects every minute could be slow

**Impact:** Increased CPU/memory usage on low-power hardware

**Mitigation:**

- Optimize query: only fetch projects with non-null cron schedules
- Consider caching parsed cron entries
- Monitor performance, optimize if needed (not premature)

### Trade-off: Simplicity vs. Features

**Decision:** Prioritize simplicity over advanced features

**What We Get:**

- Fast implementation
- Easy to understand and maintain
- Low resource usage

**What We Give Up:**

- No timezone support (yet)
- No complex scheduling patterns
- No run history for missed schedules

**Rationale:** Aligns with MVP philosophy - add advanced features based on user feedback

## Migration Plan

### Database Migration

```sql
-- Up migration
ALTER TABLE projects ADD COLUMN cron_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE projects ADD COLUMN cron_schedule TEXT NOT NULL DEFAULT '';

-- Down migration (if needed)
ALTER TABLE projects DROP COLUMN cron_enabled;
ALTER TABLE projects DROP COLUMN cron_schedule;
```

GORM will handle this automatically via `AutoMigrate()`.

### Deployment Steps

1. **Deploy backend:**
   - New `Project` entity with `CronEnabled` and `CronSchedule` fields
   - Add `go-co-op/gocron` dependency to `go.mod`
   - Create event queue adapter in `internal/adapter/scheduler`
   - Database migration runs on startup
   - Scheduler service starts alongside worker

2. **Deploy frontend:**
   - Update ProjectEditor to include cron toggle and cron input
   - Add validation for cron expressions
   - Update TypeScript interfaces

3. **No data migration needed:**
   - Existing projects get `CronEnabled=false` and `CronSchedule=''`
   - No behavior change for existing projects (scheduling disabled by default)

### Rollback Strategy

1. **Code rollback:**
   - Revert to previous version
   - Database column remains (harmless if unused)

2. **If column must be removed:**
   - Run down migration manually
   - Requires downtime

3. **Frontend rollback:**
   - Simply revert to previous version
   - No data loss (fields default to empty/disabled)

## Open Questions

1. **Cron Expression Validation:**
   - Validate on save or allow invalid expressions (fail silently)?
   - **Tentative:** Validate on save, return error if invalid

2. **UI Component:**
   - Raw text input or cron builder UI?
   - **Tentative:** Text input with validation for MVP, cron builder later

3. **Logging:**
   - Log when scheduled runs are created?
   - **Tentative:** Yes, info-level log with project ID and schedule

4. **Metrics:**
   - Track number of scheduled vs manual runs?
   - **Tentative:** Future enhancement in Stage 2 (monitoring)
