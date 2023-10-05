cidfile := $(shell mktemp -u)
cid = $$(cat $(cidfile))

all: prototype fixed.Dockerfile
	docker build -t deploy:prototype_fixed - <fixed.Dockerfile

prototype: base dummy expect_install install.sh
	docker run --cidfile $(cidfile) -v $(shell realpath .):/tmp/installer -w /tmp/installer deploy:prototype_base /tmp/installer/expect_install
	docker commit --message "install.sh" $(cid) deploy:prototype
	docker rm $(cid)
	rm $(cidfile)

base: base.Dockerfile
	docker build -t deploy:prototype_base - <base.Dockerfile

dummy:
	dd if=/dev/urandom of=dummy bs=1M count=623
