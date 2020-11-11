# rpmac
A rpm based package manager for Mac

HashiCorp Engineering Services Hackathon 2020

## Testing

### Run test repo

```bash
make test-repo
```

This will set up a test repo on a local webserver.  The repo file will be available at `localhost/test.repo`. Once that's added, the local repo should be visible to rpmac

## Functionality

### Repo management
- [x] Download metadata files from remote repos
- [x] Read information from metadata files
- [x] Store package information in local cache

### Package management
- [x] Install rpm file by name from cached metadata
- [ ] Resolve package dependencies
- [ ] Reconcile dependency versions and dependency graph
- [x] Track installed packages
- [ ] Process install/uninstall hooks

### General
- [ ] Lock during operations

### Functional Commands
- [x] Help
- [ ] Add Repo
- [x] Install
- [x] Uninstall
- [x] Search
- [ ] Upgrade
