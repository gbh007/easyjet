## MODIFIED Requirements

### Requirement: Project Entity Schema

The Project entity SHALL include cron scheduling fields for automated scheduling with enable/disable control using plain string storage.

**Current:**

```go
type Project struct {
    ID        uint      `json:"id" gorm:"column:id;not null;primarykey"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;<-:create;autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`

    Dir       string `json:"dir" gorm:"column:dir;not null"`
    GitURL    string `json:"git_url" gorm:"column:git_url;not null"`
    GitBranch string `json:"git_branch" gorm:"column:git_branch;not null"`
    Name      string `json:"name" gorm:"column:name;not null"`

    Stages []ProjectStage `json:"stages" gorm:"foreignKey:ProjectID" validate:"min=1"`
}
```

**Modified:**

```go
type Project struct {
    ID           uint      `json:"id" gorm:"column:id;not null;primarykey"`
    CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;not null;<-:create;autoCreateTime"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`
    CronEnabled  bool      `json:"cron_enabled" gorm:"column:cron_enabled;not null;default:false"`
    CronSchedule string    `json:"cron_schedule" gorm:"column:cron_schedule;type:text;not null;default:''"`

    Dir       string `json:"dir" gorm:"column:dir;not null"`
    GitURL    string `json:"git_url" gorm:"column:git_url;not null"`
    GitBranch string `json:"git_branch" gorm:"column:git_branch;not null"`
    Name      string `json:"name" gorm:"column:name;not null"`

    Stages []ProjectStage `json:"stages" gorm:"foreignKey:ProjectID" validate:"min=1"`
}
```

#### Scenario: Create project with cron scheduling

- **WHEN** a project is created with `CronEnabled=true` and valid `CronSchedule`
- **THEN** both fields are stored in the database and scheduling is active

#### Scenario: Create project without cron scheduling

- **WHEN** a project is created without cron fields
- **THEN** `CronEnabled=false` (default) and `CronSchedule=""` (empty string)

#### Scenario: Enable existing project's schedule

- **WHEN** an existing project with `CronSchedule` set is updated with `CronEnabled=true`
- **THEN** scheduling becomes active for this project

#### Scenario: Disable project schedule

- **WHEN** a project's `CronEnabled` is set to `false`
- **THEN** scheduling is paused but `CronSchedule` value is preserved

### Requirement: Project API Request/Response

The Project API SHALL accept and return the cron scheduling fields in project requests and responses.

**Current API Model:**

```json
{
  "id": 1,
  "name": "My Project",
  "dir": "/path/to/project",
  "git_url": "https://github.com/user/repo.git",
  "git_branch": "main",
  "stages": [...]
}
```

**Modified API Model:**

```json
{
  "id": 1,
  "name": "My Project",
  "dir": "/path/to/project",
  "git_url": "https://github.com/user/repo.git",
  "git_branch": "main",
  "cron_enabled": true,
  "cron_schedule": "0 5 * * *",
  "stages": [...]
}
```

#### Scenario: API returns cron scheduling fields

- **WHEN** client fetches a project with cron scheduling configured
- **THEN** the response includes both `cron_enabled` and `cron_schedule` fields

#### Scenario: API accepts cron fields in create

- **WHEN** client creates a project with `cron_enabled` and/or `cron_schedule`
- **THEN** the fields are saved and returned in the response

#### Scenario: API accepts cron fields in update

- **WHEN** client updates a project's cron fields
- **THEN** the fields are updated independently (can toggle without changing schedule)

#### Scenario: API handles empty/omitted cron schedule

- **WHEN** client creates or updates a project without `cron_schedule`
- **THEN** the field is set to empty string `""` (scheduling inactive)

### Requirement: Project Validation

The Project validation rules SHALL validate cron scheduling fields appropriately.

#### Scenario: Optional fields

- **WHEN** a project is created or updated without cron fields
- **THEN** validation passes (fields are optional with defaults)

#### Scenario: Invalid cron expression

- **WHEN** a project is created or updated with an invalid cron expression
- **THEN** validation fails and returns an error

#### Scenario: Schedule without enable flag

- **WHEN** `CronSchedule` is set but `CronEnabled=false`
- **THEN** validation passes (schedule is saved but not active)

#### Scenario: Empty string is valid

- **WHEN** `CronSchedule=""` (empty string)
- **THEN** validation passes (means "no schedule")
