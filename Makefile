.PHONY: docker shell clean release test

CURRDIR=$(shell pwd)
BUILDDIR=$(CURRDIR)/build
PKGDIR=/go/src/github.com/dcos

all: $(BUILDDIR)/dcos-checks

$(BUILDDIR)/dcos-checks: docker $(shell find "$(CURRDIR)" -name "*.go")
	mkdir -p $(BUILDDIR)
	docker run -v $(BUILDDIR):$(BUILDDIR) \
			   -v $(CURRDIR):$(PKGDIR)/dcos-checks \
			   -w $(PKGDIR)/dcos-checks \
			   --rm \
			   dcos/dcos-checks-test \
			   bash -c "go build -o $(@) ."

test: $(BUILDDIR)/dcos-checks
	docker run -v $(CURRDIR):$(PKGDIR)/dcos-checks \
			   -w $(PKGDIR)/dcos-checks \
			   --rm \
			   --privileged \
			   dcos/dcos-checks-test \
			   bash -c "./script/test.sh unit"

docker:
	docker build --rm --force-rm -f Dockerfile -t dcos/dcos-checks-test .

clean:
	rm -rf $(BUILDDIR)

shell:
	docker run --rm -it \
		--privileged \
		-v $(CURRDIR):$(PKGDIR)/dcos-checks \
		-w $(PKGDIR)/dcos-checks \
		dcos/dcos-checks-test /bin/bash

build: all
