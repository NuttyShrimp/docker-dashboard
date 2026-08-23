package jobs

import (
	"context"
	"time"

	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
	"go.uber.org/zap"
)

type Jobs struct {
	jobs   []Job
	logger *zap.Logger
}

type Job interface {
	Name() string
	Run(ctx context.Context) error
}

func NewJobs(service *service.Service) *Jobs {
	logger := zap.L().Named("jobs")
	jobs := &Jobs{
		logger: logger,
		jobs: []Job{
			NewDockerJob(service, logger),
		},
	}

	return jobs
}

func (j *Jobs) process(ctx context.Context) {
	j.logger.Info("Running jobs")
	for _, job := range j.jobs {
		if err := job.Run(ctx); err != nil {
			j.logger.Error("job failed", zap.String("job", "docker"), zap.Error(err))
		}
	}
}

func (j *Jobs) Run(ctx context.Context) {
	tick := time.NewTicker(time.Second * 30)
	j.process(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			j.process(ctx)
		}
	}
}
