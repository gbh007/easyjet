package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/gbh007/easyjet/internal/core/port"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler  gocron.Scheduler
	eventQueue <-chan entity.Event
	service    port.Service
	logger     *slog.Logger
	jobs       map[uint]gocron.Job
	jobsMu     *sync.Mutex
}

func NewScheduler(logger *slog.Logger, ps port.PubSub, service port.Service) *Scheduler {
	return &Scheduler{
		eventQueue: ps.SubscribeEvent("scheduler", 100),
		service:    service,
		logger:     logger,
		jobs:       make(map[uint]gocron.Job),
		jobsMu:     &sync.Mutex{},
	}
}

func (s *Scheduler) start(ctx context.Context) error {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("create scheduler: %w", err)
	}

	s.scheduler = sched

	go s.processEvents(ctx)

	if err := s.loadExistingJobs(ctx); err != nil {
		return fmt.Errorf("load existing jobs: %w", err)
	}

	sched.Start()

	return nil
}

func (s *Scheduler) Serve(ctx context.Context) error {
	err := s.start(ctx)
	if err != nil {
		return fmt.Errorf("start scheduler: %w", err)
	}

	<-ctx.Done()

	err = s.scheduler.Shutdown()
	if err != nil {
		return fmt.Errorf("stop scheduler: %w", err)
	}

	return nil
}

func (s *Scheduler) loadExistingJobs(ctx context.Context) error {
	projects, err := s.service.Projects(ctx)
	if err != nil {
		return fmt.Errorf("fetch projects: %w", err)
	}

	for _, project := range projects {
		if project.CronEnabled && project.CronSchedule != "" {
			err := s.registerJob(ctx, project.ID)
			if err != nil {
				return fmt.Errorf("register project %d: %w", project.ID, err)
			}
		}
	}

	return nil
}

func (s *Scheduler) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-s.eventQueue:
			if !ok {
				return
			}
			err := s.handleEvent(ctx, event)
			if err != nil {
				s.logger.Error("fail handle pubsub event", "error", err)
			}
		}
	}
}

func (s *Scheduler) handleEvent(ctx context.Context, event entity.Event) (err error) {
	switch event.Type { //nolint:exhaustive // не требуется поддержка всех вариантов
	case entity.EventProjectCreated:
		err = s.registerJob(ctx, event.ProjectID)
		if err != nil {
			return fmt.Errorf("register job: %w", err)
		}
	case entity.EventProjectUpdated:
		err = s.removeJob(event.ProjectID)
		if err != nil {
			return fmt.Errorf("remove job: %w", err)
		}

		err = s.registerJob(ctx, event.ProjectID)
		if err != nil {
			return fmt.Errorf("register job: %w", err)
		}
	case entity.EventProjectDeleted:
		err = s.removeJob(event.ProjectID)
		if err != nil {
			return fmt.Errorf("remove job: %w", err)
		}
	}

	return nil
}

func (s *Scheduler) registerJob(ctx context.Context, projectID uint) error {
	project, err := s.service.Project(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}

	if !project.CronEnabled || project.CronSchedule == "" {
		return nil
	}

	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	job, err := s.scheduler.NewJob(
		gocron.CronJob(project.CronSchedule, false),
		gocron.NewTask(func(ctx context.Context) error {
			_, err := s.service.RunProject(ctx, projectID)
			return err
		}),
		gocron.WithName(fmt.Sprintf("project-%d", projectID)),
	)
	if err != nil {
		return fmt.Errorf("create cron job: %w", err)
	}

	s.jobs[projectID] = job

	return nil
}

func (s *Scheduler) removeJob(projectID uint) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	job, exists := s.jobs[projectID]

	if !exists {
		return nil
	}

	return s.scheduler.RemoveJob(job.ID())
}
