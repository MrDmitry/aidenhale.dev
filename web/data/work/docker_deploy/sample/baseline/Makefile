all: dummy Dockerfile expect_install install.sh
	docker build -t deploy:baseline .

dummy:
	dd if=/dev/urandom of=dummy bs=1M count=623
