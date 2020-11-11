# This should assist testing and stuff

.PHONY: test-packages
test-packages: # test-package creates an rpm package for testing installation
	cd test/package/rpmac-test && fpm -p ../../ -n rpmac-test -v 1.0 -s dir -t rpm -a all --prefix '/' --verbose .
	cd test/package/kpfoo && fpm -p ../../ -n kpfoo -v 1.1.0 -s dir -t rpm -a all --prefix '/' --verbose .
	cd test/package/kpbar && fpm -p ../../ -n kpbar -v 0.5.1 -s dir -t rpm -a all --prefix '/' -d 'kpfoo' --verbose .

.PHONY: test-repo
test-repo: # test-repo creates a local test repository for testing package management
	docker stop test-repo & docker rm test-repo & true
	docker run --name test-repo -v ${PWD}/test/repo:/usr/share/nginx/html:ro -p 80:80 -d nginx
