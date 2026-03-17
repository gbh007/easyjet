package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gbh007/easyjet/internal/adapter/scheduler"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/gbh007/easyjet/internal/core/port"
	"github.com/go-co-op/gocron/v2"
)

// Scheduler manages cron-based scheduling for projects.
type Scheduler struct {
	scheduler  gocron.Scheduler
	eventQueue *scheduler.EventQueue
	db         port.Database
	service    port.Service
	logger     *slog.Logger
	jobs       map[uint]gocron.Job // project ID -> job mapping
	jobsMu     sync.RWMutex        // protect jobs map
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(logger *slog.Logger, eventQueue *scheduler.EventQueue, db port.Database, service port.Service) *Scheduler {
	return &Scheduler{
		eventQueue: eventQueue,
		db:         db,
		service:    service,
		logger:     logger,
		jobs:       make(map[uint]gocron.Job),
	}
}

// Start initializes and starts the scheduler.
func (s *Scheduler) Start(ctx context.Context) error {
	// Initialize gocron scheduler
	sched, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("create scheduler: %w", err)
	}

	s.scheduler = sched

	// Start event queue consumer
	go s.processEvents(ctx)

	// Load existing cron jobs from database
	if err := s.loadExistingJobs(ctx); err != nil {
		return fmt.Errorf("load existing jobs: %w", err)
	}

	s.logger.Info("scheduler started")
	return nil
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop(ctx context.Context) error {
	if s.scheduler != nil {
		if err := s.scheduler.Shutdown(); err != nil {
			return fmt.Errorf("shutdown scheduler: %w", err)
		}
	}

	if s.eventQueue != nil {
		s.eventQueue.Close()
	}

	s.logger.Info("scheduler stopped")
	return nil
}

// loadExistingJobs loads all enabled cron jobs from the database on startup.
func (s *Scheduler) loadExistingJobs(ctx context.Context) error {
	projects, err := s.db.Projects(ctx)
	if err != nil {
		return fmt.Errorf("fetch projects: %w", err)
	}

	for _, project := range projects {
		if project.CronEnabled && project.CronSchedule != "" {
			if err := s.registerJob(ctx, project); err != nil {
				s.logger.Error("failed to register job on startup",
					"project_id", project.ID,
					"error", err)
			} else {
				s.logger.Debug("registered job on startup",
					"project_id", project.ID,
					"schedule", project.CronSchedule)
			}
		}
	}

	return nil
}

// processEvents consumes events from the event queue and processes them.
func (s *Scheduler) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-s.eventQueue.Subscribe():
			if !ok {
				return
			}
			s.handleEvent(ctx, event)
		}
	}
}

// handleEvent processes a single scheduler event.
// FIXME(ai-shit): много бесполезного мусора функцию надо сильно упростить.
func (s *Scheduler) handleEvent(ctx context.Context, event scheduler.SchedulerEvent) {
	s.logger.Debug("processing scheduler event",
		"event_type", event.Type,
		"project_id", event.ProjectID)

	switch event.Type {
	case scheduler.EventCreated:
		s.handleProjectCreated(ctx, event)
	case scheduler.EventUpdated:
		s.handleProjectUpdated(ctx, event)
	case scheduler.EventDeleted:
		s.handleProjectDeleted(ctx, event)
	}
}

// handleProjectCreated registers a new cron job for a created project.
func (s *Scheduler) handleProjectCreated(ctx context.Context, event scheduler.SchedulerEvent) {
	if !event.Enabled || event.Schedule == "" {
		s.logger.Debug("skipping job registration for project",
			"project_id", event.ProjectID,
			"enabled", event.Enabled,
			"schedule", event.Schedule)
		return
	}

	// Fetch full project details
	project, err := s.db.Project(ctx, event.ProjectID)
	if err != nil {
		s.logger.Error("failed to fetch project for job registration",
			"project_id", event.ProjectID,
			"error", err)
		return
	}

	if err := s.registerJob(ctx, project); err != nil {
		s.logger.Error("failed to register job for created project",
			"project_id", event.ProjectID,
			"error", err)
		return
	}

	s.logger.Info("registered cron job for new project",
		"project_id", project.ID,
		"schedule", project.CronSchedule)
}

// handleProjectUpdated updates or removes a cron job when project settings change.
func (s *Scheduler) handleProjectUpdated(ctx context.Context, event scheduler.SchedulerEvent) {
	// First remove existing job if any
	s.removeJob(event.ProjectID)

	// If scheduling is enabled, register new job
	if event.Enabled && event.Schedule != "" {
		project, err := s.db.Project(ctx, event.ProjectID)
		if err != nil {
			s.logger.Error("failed to fetch project for job update",
				"project_id", event.ProjectID,
				"error", err)
			return
		}

		if err := s.registerJob(ctx, project); err != nil {
			s.logger.Error("failed to update job for project",
				"project_id", event.ProjectID,
				"error", err)
			return
		}

		s.logger.Info("updated cron job for project",
			"project_id", project.ID,
			"schedule", project.CronSchedule)
	} else {
		s.logger.Info("removed cron job for project (disabled)",
			"project_id", event.ProjectID)
	}
}

// handleProjectDeleted removes a cron job when a project is deleted.
func (s *Scheduler) handleProjectDeleted(_ context.Context, event scheduler.SchedulerEvent) {
	s.removeJob(event.ProjectID)
	s.logger.Info("removed cron job for deleted project",
		"project_id", event.ProjectID)
}

// registerJob registers a cron job for a project.
func (s *Scheduler) registerJob(ctx context.Context, project entity.Project) error {
	if !project.CronEnabled || project.CronSchedule == "" {
		return nil
	}

	// Create cron job using gocron.Cron() method
	job, err := s.scheduler.NewJob(
		gocron.CronJob(project.CronSchedule, false), // false = don't use seconds
		// FIXME(ai-shit): проверить про передачу контекста, чтобы получать правильный а не корневой.
		gocron.NewTask(s.executeScheduledRun, ctx, project.ID),
		gocron.WithName(fmt.Sprintf("project-%d", project.ID)),
	)
	if err != nil {
		return fmt.Errorf("create cron job: %w", err)
	}

	s.jobsMu.Lock()
	s.jobs[project.ID] = job
	s.jobsMu.Unlock()

	return nil
}

// removeJob removes a cron job for a project.
func (s *Scheduler) removeJob(projectID uint) {
	s.jobsMu.RLock()
	job, exists := s.jobs[projectID]
	s.jobsMu.RUnlock()

	if !exists {
		return
	}

	if err := s.scheduler.RemoveJob(job.ID()); err != nil {
		s.logger.Error("failed to remove job",
			"project_id", projectID,
			"error", err)
		return
	}

	s.jobsMu.Lock()
	delete(s.jobs, projectID)
	s.jobsMu.Unlock()
}

// executeScheduledRun creates a pending run when a cron schedule triggers.
func (s *Scheduler) executeScheduledRun(ctx context.Context, projectID uint) {
	s.logger.Info("cron schedule triggered, creating pending run",
		"project_id", projectID)

	// Create a pending run
	runID, err := s.service.RunProject(ctx, projectID)
	if err != nil {
		s.logger.Error("failed to create pending run for scheduled execution",
			"project_id", projectID,
			"error", err)
		return
	}

	s.logger.Info("created pending run for scheduled execution",
		"project_id", projectID,
		"run_id", runID)
}
