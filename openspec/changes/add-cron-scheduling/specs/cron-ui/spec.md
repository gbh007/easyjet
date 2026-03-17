## ADDED Requirements

### Requirement: Cron Schedule Toggle Control

The Project Editor SHALL include a toggle switch to enable/disable cron scheduling.

#### Scenario: Display toggle control

- **WHEN** user opens the project editor (create or edit mode)
- **THEN** a toggle/switch control for enabling/disabling scheduling is displayed

#### Scenario: Toggle reflects current state

- **WHEN** editing a project with `CronEnabled=true`
- **THEN** the toggle is in the "on" position

#### Scenario: Toggle off preserves schedule

- **WHEN** user toggles scheduling off
- **THEN** the `CronSchedule` value is preserved but scheduling is disabled

### Requirement: Cron Schedule Input Field

The Project Editor SHALL include a cron schedule input field for configuring automated scheduling.

#### Scenario: Display cron schedule field

- **WHEN** user opens the project editor
- **THEN** a cron schedule input field is displayed below Git configuration fields

#### Scenario: Edit existing cron schedule

- **WHEN** user edits a project with an existing cron schedule
- **THEN** the field is pre-populated with the current cron expression

#### Scenario: Empty cron schedule

- **WHEN** a project has no cron schedule (NULL)
- **THEN** the input field is empty but editable

#### Scenario: Input disabled when toggle is off

- **WHEN** the cron toggle is in the "off" position
- **THEN** the cron schedule input may be visually disabled or greyed out (optional UX enhancement)

### Requirement: Cron Expression Validation

The frontend SHALL validate cron expressions before submitting to the API.

#### Scenario: Valid cron expression

- **WHEN** user enters a valid 5-field cron expression
- **THEN** the form accepts the value and allows submission

#### Scenario: Invalid cron expression

- **WHEN** user enters an invalid cron expression
- **THEN** the form displays a validation error and prevents submission

#### Scenario: Empty cron schedule

- **WHEN** user leaves the cron schedule field empty
- **THEN** the field is submitted as NULL (no scheduling)

### Requirement: Cron Schedule Submission

The frontend SHALL include cron scheduling fields in project create/update API requests.

#### Scenario: Create project with cron scheduling

- **WHEN** user creates a new project with cron scheduling enabled
- **THEN** both `cron_enabled` and `cron_schedule` fields are included in the POST request

#### Scenario: Update project cron schedule

- **WHEN** user updates an existing project's cron schedule
- **THEN** both `cron_enabled` and `cron_schedule` fields are included in the PUT request

#### Scenario: Toggle scheduling off

- **WHEN** user disables the cron toggle and saves
- **THEN** `cron_enabled=false` is submitted (schedule value preserved on server)

#### Scenario: Toggle scheduling on

- **WHEN** user enables the cron toggle with a valid schedule and saves
- **THEN** `cron_enabled=true` and `cron_schedule` are submitted

### Requirement: Cron Schedule Display

The Project Details page SHALL display the current cron schedule status if configured.

#### Scenario: Display active schedule

- **WHEN** a project has `CronEnabled=true` and a cron schedule
- **THEN** the schedule is displayed with an indicator showing it's active (e.g., "Cron: 0 5 \* \* \* ✓")

#### Scenario: Display disabled schedule

- **WHEN** a project has `CronEnabled=false` but has a saved cron schedule
- **THEN** the schedule is displayed with an indicator showing it's paused (e.g., "Cron: 0 5 \* \* \* ⏸")

#### Scenario: No schedule

- **WHEN** a project has no cron schedule
- **THEN** no cron information is displayed or "Manual only" is shown

### Requirement: User Experience

The cron schedule UI SHALL provide helpful guidance for users.

#### Scenario: Cron format hint

- **WHEN** user focuses the cron schedule input
- **THEN** a hint showing the cron format (minute hour day month weekday) is displayed

#### Scenario: Example cron expressions

- **WHEN** user views the cron schedule field
- **THEN** common examples are visible (e.g., "0 5 \* \* \* = daily at 5:00 AM")

#### Scenario: Responsive design

- **WHEN** user accesses the form on different screen sizes
- **THEN** the cron toggle and input remain usable and properly formatted
