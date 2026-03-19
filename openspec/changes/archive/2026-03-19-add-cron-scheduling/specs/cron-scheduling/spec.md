## ADDED Requirements

### Requirement: Scheduler Service Initialization

The system SHALL initialize a cron scheduler service on application startup that runs alongside the existing worker service.

#### Scenario: Scheduler starts with application

- **WHEN** EasyJet application starts
- **THEN** scheduler service starts automatically as a background goroutine

#### Scenario: Scheduler stops gracefully

- **WHEN** application receives shutdown signal (SIGINT, SIGTERM, SIGQUIT)
- **THEN** scheduler stops gracefully via context cancellation

### Requirement: Cron Expression Parsing

The system SHALL parse and validate standard cron expressions using the `go-co-op/gocron` library.

#### Scenario: Valid cron expression accepted

- **WHEN** a valid 5-field cron expression is provided (e.g., `0 5 * * *`)
- **THEN** the expression is parsed successfully and scheduled

#### Scenario: Invalid cron expression rejected

- **WHEN** an invalid cron expression is provided
- **THEN** the system returns a validation error and does not schedule the run

### Requirement: Cron Schedule Storage

The system SHALL store cron schedules as plain strings (not pointers) with empty string meaning "no schedule".

#### Scenario: Store valid cron schedule

- **WHEN** a project has a cron schedule
- **THEN** the schedule is stored as a plain string in the database

#### Scenario: No schedule represented

- **WHEN** a project has no cron schedule
- **THEN** the field is an empty string `""` (not NULL)

### Requirement: Cron Schedule Enable/Disable Toggle

The system SHALL respect the `CronEnabled` field to determine if scheduling is active for a project.

#### Scenario: Enabled project with schedule

- **WHEN** `CronEnabled=true` AND `CronSchedule` is not empty
- **THEN** the scheduler evaluates and triggers runs for this project

#### Scenario: Disabled project with schedule

- **WHEN** `CronEnabled=false` AND `CronSchedule` is set
- **THEN** the scheduler does NOT trigger runs (schedule is paused)

#### Scenario: Project without schedule

- **WHEN** `CronSchedule` is empty string (regardless of `CronEnabled`)
- **THEN** the scheduler does NOT trigger runs for this project

### Requirement: Event Queue Adapter for Configuration Changes

The system SHALL use an event queue adapter with channels to handle cron configuration changes dynamically.

#### Scenario: Project created with cron schedule

- **WHEN** a new project is created with cron scheduling enabled
- **THEN** a "project created" event is sent to the event queue
- **AND** the scheduler picks up the event and registers the cron job

#### Scenario: Project cron settings updated

- **WHEN** an existing project's cron settings are modified
- **THEN** a "project updated" event is sent to the event queue
- **AND** the scheduler picks up the event and updates/removes the cron job accordingly

#### Scenario: Project deleted

- **WHEN** a project with cron scheduling is deleted
- **THEN** a "project deleted" event is sent to the event queue
- **AND** the scheduler picks up the event and removes the cron job

#### Scenario: Event queue processing

- **WHEN** events are in the queue
- **THEN** the scheduler processes them sequentially via channel
- **AND** updates its internal job registry accordingly

### Requirement: Scheduled Run Creation

The system SHALL create a pending project run when a cron schedule triggers.

#### Scenario: Cron schedule triggers

- **WHEN** the current time matches a project's cron schedule
- **AND** `CronEnabled=true`
- **THEN** a new pending run is created in the database

### Requirement: Scheduler Integration with Worker

The scheduler SHALL integrate with the existing worker infrastructure by creating pending runs.

#### Scenario: Worker processes scheduled runs

- **WHEN** scheduler creates a pending run
- **THEN** the existing worker picks it up and executes it normally

#### Scenario: No distinction between manual and scheduled runs

- **WHEN** a run is executed
- **THEN** the execution logic is identical regardless of whether it was triggered manually or by schedule

### Requirement: Scheduler Logging

The system SHALL log scheduler activities for observability.

#### Scenario: Log scheduled run creation

- **WHEN** a scheduled run is created
- **THEN** an INFO-level log is written with project ID and cron expression

#### Scenario: Log scheduler errors

- **WHEN** an error occurs during scheduling (e.g., database error)
- **THEN** an ERROR-level log is written with error details

#### Scenario: Log event queue processing

- **WHEN** an event is processed from the queue
- **THEN** a DEBUG-level log is written with event type and project ID

### Requirement: Single Instance Execution

The system SHALL operate as a single-instance scheduler with single concurrent worker.

#### Scenario: Single scheduler instance

- **WHEN** EasyJet runs
- **THEN** only one scheduler instance exists
- **AND** no distributed coordination is needed

#### Scenario: No concurrent run protection needed

- **WHEN** a cron schedule triggers
- **THEN** the run is created without checking for existing runs
- **BECAUSE** single worker processes runs sequentially
