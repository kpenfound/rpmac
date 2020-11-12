# rpmac
A rpm based package manager for Mac

HashiCorp Engineering Services Hackathon 2020


## Functionality

### Repo management
- [x] Download metadata files from remote repos
- [x] Read information from metadata files
- [x] Store package information in local cache

### Package management
- [x] Install rpm file by name from cached metadata
- [x] Track installed packages
- [x] Install a specific version of a package
- [ ] Reconcile dependency versions and dependency graph
- [ ] Process install/uninstall hooks

### General
- [x] Lock during operations

### Functional Commands
- [x] Help
- [x] Add Repo
- [x] Install
- [x] Uninstall
- [x] Search
- [x] Upgrade


## Testing

```bash
go build
```


### Run test repo

```bash
make test-repo
```

This will set up a test repo on a local webserver.  The repo file will be available at `http://localhost/test.repo`. Once that's added, the local repo should be visible to rpmac

### Add repo

```bash
./rpmac add-repo ${PWD}/test/repo/test.repo
```

### Search for a package

```bash
./rpmac search kpfoo
```

### Install a package

```bash
./rpmac install kpfoo
```

### Uninstall a package

```bash
./rpmac uninstall kpfoo
```

### Upgrade a package

```bash
./rpmac install kpfoo=1.1.0
./rpmac upgrade kpfoo
```