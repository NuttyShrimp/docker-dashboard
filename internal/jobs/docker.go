package jobs

import (
	"context"

	"github.com/moby/moby/client"
	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
	"go.uber.org/zap"
)

type DockerJob struct {
	// service  *service.Service
	logger    *zap.Logger
	subdomain *service.Subdomain
	resource  *service.Resource
}

func NewDockerJob(service *service.Service, logger *zap.Logger) *DockerJob {
	job := DockerJob{
		subdomain: service.NewSubdomain(),
		resource:  service.NewResource(),
		logger:    logger.Named("docker"),
	}

	return &job
}

func (d *DockerJob) Name() string {
	return "docker"
}

func (d *DockerJob) Run(ctx context.Context) error {
	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("docker-dashboard/1.0.0"),
	)
	if err != nil {
		return err
	}
	// nolint:errcheck // we don't care if the client is closed correctly
	defer apiClient.Close()

	result, err := apiClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for i := range result.Items {
		container := result.Items[i]
		exposedPorts := []int{}
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				exposedPorts = append(exposedPorts, int(port.PublicPort))
			}
		}
		// We don't need unmapped resources
		if len(exposedPorts) == 0 {
			continue
		}
		if len(container.Names) == 0 {
			zap.L().Warn("container without name", zap.String("id", container.ID))
		}
		project := container.Labels["com.docker.compose.project"]

		// Check for subdomain based on custom label, then project label, then unowned subdomain
		// Check if subdomain exists & create when missing
		if _, ok := d.subdomain.Get(project); !ok {
			_ = d.subdomain.Create(project)
		}

		// Check that resource exists
		sanitizedName := d.resource.SanitizeCtName(container.Names[0], project)
		resource, ok := d.resource.Get(sanitizedName, project)
		if !ok {
			// Create when missing
			resource, err = d.resource.Create(sanitizedName, project)
			if err != nil {
				zap.L().Error("Failed to create resource", zap.String("resource", sanitizedName), zap.String("subdomain", project))
			}
		}

		resource.Ports = append(resource.Ports, exposedPorts...)
	}

	return nil
}
