package repository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kpenfound/rpmac/constants"
	"github.com/kpenfound/rpmac/rpm"
	"github.com/kpenfound/rpmac/util"
)

// Repositories struct for all local repositories
type Repositories struct {
	Repositories []*Repository
}

// InitRepositories initializes the repository objects
func InitRepositories() (Repositories, error) {
	r := Repositories{}
	repos := []*Repository{}
	repodir := util.ReplaceHome(constants.RepoDir)
	files, err := ioutil.ReadDir(repodir)
	if err != nil {
		return r, err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".repo") {
			fname := filepath.Join(repodir, f.Name())
			repo, err := ReadRepoFile(fname)
			if err != nil {
				return r, err
			}
			for _, rs := range repo { // A repo file can contain multiple repos
				repos = append(repos, &rs)
			}
		}

	}

	r.Repositories = repos
	err = r.Sync()
	if err != nil {
		return r, err
	}
	err = r.LoadCache()
	if err != nil {
		return r, err
	}

	return r, nil
}

// Sync repository metadata with local cache
func (r *Repositories) Sync() error {
	for _, repo := range r.Repositories {
		fmt.Printf("Syncing %s\n", repo.Name)
		if repo.Enabled {
			err := repo.Sync()
			if err != nil {
				return err
			}
		}
		fmt.Printf("Cache files: %+v\n", repo.CacheFiles)
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

// LoadCache load package cache of all enabled repos
func (r *Repositories) LoadCache() error {
	for _, repo := range r.Repositories {
		fmt.Printf("Loading cache for %s\n", repo.Name)
		fmt.Printf("Cache files: %+v\n", repo.CacheFiles)
		if repo.Enabled {
			err := repo.LoadCache()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
