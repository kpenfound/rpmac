# rpmac
A rpm based package manager for Mac

## Testing

### Run test repo

```bash
make test-repo
```

This will set up a test repo on a local webserver.  The repo file will be available at `localhost/test.repo`. Once that's added, the local repo should be visible to rpmac

## Functionality

### Repo management
- [ ] Download metadata files from remote repos
- [ ] Read information from metadata files
- [ ] Store package information in local cache

### Package management
- [ ] Install rpm file by name from cached metadata
- [ ] Resolve package dependencies
- [ ] Reconcile dependency versions and dependency graph
- [ ] Track installed packages
- [ ] Process install/uninstall hooks

### General
- [ ] Lock during operations

### Functional Commands
- [x] Help
- [ ] Add Repo
- [ ] Install
- [ ] Uninstall
- [ ] Search
- [ ] Upgrade
