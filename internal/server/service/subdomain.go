package service

import (
	"github.com/nuttyshrimp/docker-dashboard/internal/database/model"
	"github.com/nuttyshrimp/docker-dashboard/internal/database/repository"
	"github.com/nuttyshrimp/docker-dashboard/pkg/utils"
)

type Subdomain struct {
	subdomain *repository.Subdomain
}

func (s *Service) NewSubdomain() *Subdomain {
	return &Subdomain{
		subdomain: s.repo.NewSubdomain(),
	}
}

func (s *Subdomain) GetAll() []model.Subdomain {
	subs := utils.SliceMap(s.subdomain.GetAll(), func(subdomain *model.Subdomain) model.Subdomain {
		return *subdomain
	})
	return subs
}

func (s *Subdomain) Get(name string) (model.Subdomain, bool) {
	sub, _ := s.subdomain.GetByName(name)
	if sub == nil {
		return model.Subdomain{}, false
	}
	return *sub, true
}

func (s *Subdomain) Create(name string) model.Subdomain {
	return *s.subdomain.Create(name)
}
