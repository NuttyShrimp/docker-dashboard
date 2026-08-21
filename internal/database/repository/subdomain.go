package repository

import (
	"github.com/nuttyshrimp/docker-dashboard/internal/database/model"
	"github.com/nuttyshrimp/docker-dashboard/pkg/utils"
)

type Subdomain struct {
	repo *Repository
}

func (r *Repository) NewSubdomain() *Subdomain {
	model := &Subdomain{repo: r}
	return model
}

func (s *Subdomain) GetAll() []*model.Subdomain {
	s.repo.lock.RLock()
	defer s.repo.lock.RUnlock()
	return s.repo.domains
}

func (s *Subdomain) GetByName(name string) (*model.Subdomain, bool) {
	s.repo.lock.RLock()
	defer s.repo.lock.RUnlock()
	domain, _ := utils.SliceFind(s.repo.domains, func(sd *model.Subdomain) bool {
		return sd.Name == name
	})
	if domain == nil {
		return nil, false
	}
	return *domain, true
}

func (s *Subdomain) Create(name string, resources ...*model.Resource) *model.Subdomain {
	s.repo.lock.Lock()
	defer s.repo.lock.Unlock()
	subdomain := &model.Subdomain{
		Name:      name,
		Resources: resources,
	}
	s.repo.domains = append(s.repo.domains, subdomain)
	return subdomain
}
