package repository

import "github.com/nuttyshrimp/docker-dashboard/internal/database/model"

type Resource struct {
	repo *Repository
}

func (r *Repository) NewResource() *Resource {
	return &Resource{repo: r}
}

func (r *Resource) GetByName(name string) *model.Resource {
	r.repo.lock.RLock()
	defer r.repo.lock.RUnlock()
	for _, domain := range r.repo.domains {
		for _, resource := range domain.Resources {
			if resource.Name == name {
				return resource
			}
		}
	}

	return nil
}
