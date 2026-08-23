package service

import (
	"strings"

	"github.com/nuttyshrimp/docker-dashboard/internal/database/model"
	"github.com/nuttyshrimp/docker-dashboard/internal/database/repository"
	"github.com/nuttyshrimp/docker-dashboard/pkg/utils"
)

type Resource struct {
	resource  *repository.Resource
	subdomain *repository.Subdomain
}

func (s *Service) NewResource() *Resource {
	return &Resource{
		resource:  s.repo.NewResource(),
		subdomain: s.repo.NewSubdomain(),
	}
}

func (r *Resource) Get(name, project string) (*model.Resource, bool) {
	// Cheap search to see if project (subdomain) with container (resource) already exists
	subdomain, ok := r.subdomain.GetByName(project)
	if ok {
		if resource, ok := utils.SliceFind(subdomain.Resources, func(r *model.Resource) bool {
			return r.Name == name
		}); ok {
			return *resource, true
		}
	}
	// Expensive, all containers in projects
	resource := r.resource.GetByName(name)
	if resource == nil {
		return nil, false
	}

	return resource, true
}

func (r *Resource) GetPortToResource() map[int]*model.Resource {
	mappedResources := make(map[int]*model.Resource)
	subdomains := r.subdomain.GetAll()

	for _, domain := range subdomains {
		for _, resource := range domain.Resources {
			for _, port := range resource.Ports {
				mappedResources[port] = resource
			}
		}
	}

	return mappedResources
}

func (r *Resource) Create(name, project string) (*model.Resource, error) {
	subdomain, ok := r.subdomain.GetByName(project)
	if !ok {
		return nil, utils.Errorf("subdomain %s is not available", project)
	}
	resource := &model.Resource{
		Name: name,
	}
	if subdomain == nil {
		subdomain = r.subdomain.Create(project, resource)
	}
	resource.Subdomain = subdomain
	subdomain.Resources = append(subdomain.Resources, resource)
	return resource, nil

}

// Removes the folder name (the project name normally) prefix from the container name when it exists
func (r *Resource) SanitizeCtName(name, project string) string {
	// name: infisical -> do not modify
	// name: infisical-dev -> do not modify
	// name: infisical-infisical-dev -> modify to : infisical-dev
	if name[0] == '/' {
		name = strings.TrimLeft(name, "/")
	}

	if strings.Contains(name, project+"-") {
		name = strings.Replace(name, project+"-", "", 1)
	}

	return name
}
