package repository

import (
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

// RepoPackage used for Query responses
type RepoPackage struct {
	Repository *Repository
	Package    *rpm.RPM
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

	return r, nil
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

// QueryOptions options for Query
type QueryOptions struct {
	// Name
	Name string
	// FuzzyVersion, like 1.0.0-beta1, 1.0.0, 1.0, or 1
	FuzzyVersion string
}

// MakeQueryOptions creates a QueryOptions struct from a package string like firefox=33.*
func MakeQueryOptions(packagestr string) QueryOptions {
	qo := QueryOptions{}

	// Separate name from version
	parts := strings.Split(packagestr, "=")
	qo.Name = parts[0]
	if len(parts) > 1 {
		qo.FuzzyVersion = parts[1]
	}

	return qo
}

// Query for a package by name in local cache
func (r *Repositories) Query(opts QueryOptions) (*RepoPackage, error) {
	p := RepoPackage{}

	for _, repo := range r.Repositories {
		rpm, err := repo.Query(opts)
		if err != nil && err != constants.ErrorPackageNotFound {
			return &p, err
		}

		if err != constants.ErrorPackageNotFound {
			p = RepoPackage{
				Repository: repo,
				Package:    rpm,
			}
			return &p, nil
		}
	}
	return &p, constants.ErrorPackageNotFound
}

// LoadCache load package cache of all enabled repos
func (r *Repositories) LoadCache() error {
	for _, repo := range r.Repositories {
		if repo.Enabled {
			err := repo.LoadCache()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
