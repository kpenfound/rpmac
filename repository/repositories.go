package repository

import (
	"github.com/kpenfound/rpmac/rpm"
)

// Repositories struct for all local repositories
type Repositories struct {
	Repositories []Repository
}

// Sync repository metadata with local cache
func (r *Repositories) Sync() error {
	for _, repo := range r.Repositories {
		if repo.Enabled {
			err := repo.Sync()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Query for a package by name in local cache
func (r *Repositories) Query(name string) (rpm.RPM, error) {
	p := rpm.RPM{}

	for _, repo := range r.Repositories {
		p, err := repo.Query(name)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
