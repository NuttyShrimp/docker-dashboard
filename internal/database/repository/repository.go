package repository

import (
	"sync"

	"github.com/nuttyshrimp/docker-dashboard/internal/database/model"
)

// Will hold our data in memory
type Repository struct {
	lock    sync.RWMutex
	domains []*model.Subdomain
}

func New() *Repository {
	return &Repository{
		lock: sync.RWMutex{},
	}
}
