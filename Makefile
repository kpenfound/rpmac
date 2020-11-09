# This should assist testing and stuff

.PHONY: test-package
test-package: # test-package creates an rpm package for testing installation
	fpm -p ./ -n rpmac-test -v 1.0 -s dir -t rpm -a all --prefix '/' --verbose ./test/package

.PHONY: test-repo
test-repo: # test-repo creates a local test repository for testing package management
	docker rm test-repo & true
	docker run --name test-repo -v ${PWD}/test/repo:/usr/share/nginx/html:ro -p 80:80 -d nginx
