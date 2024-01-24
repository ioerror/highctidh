library: _prep
	cd src; make;
	cp -v src/*.so dist/;

release: clean _prep
	cd src; make -f Makefile.packages sdist;
	./misc/docker-multi-arch-package-build.sh;
	echo "Watch the build process: tail -f docker_build_output/*/*/build.log";
	echo "Please run: cd src/ && make -f Makefile.packages pypi-upload";

packages: _prep
	cd src/; make -f Makefile.packages packages;

deb: _prep
	cd src/; make -f Makefile.packages deb;

wheel: _prep
	cd src/; make -f Makefile.packages wheel;

pytest:
	cd src; pytest-3 -v -n auto --doctest-modules;

update-golang-modules:
	cd src/; make -f Makefile.packages update-golang-modules;

install:
	cd src; make install;

test:
	cd src; make test;

test-go:
	cd src; go test -v ./...;
	cd src/ctidh511; go test -v ./...;
	cd src/ctidh512; go test -v ./...;
	cd src/ctidh1024; go test -v ./...;
	cd src/ctidh2048; go test -v ./...;

examples:
	cd src; make examples;

examples-run:
	cd src; time ./example-ctidh511;
	cd src; time ./example-ctidh512;
	cd src; time ./example-ctidh1024;
	cd src; time ./example-ctidh2048;

_prep:
	-mkdir src/dist;
	-ln -s src/build build;
	-ln -s src/dist dist;

clean:
	-rm -rf build;
	-rm -rf dist;
	-rm -rf docker_build_output;
	-cd src; make clean;
	-cd src; make -f Makefile.packages clean;
