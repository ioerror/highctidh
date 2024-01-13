library: _prep
	cd src; make;
	cp -v src/*.so dist/;

packages: _prep
	cd src/; make -f Makefile.packages packages;

deb: _prep
	cd src/; make -f Makefile.packages deb

wheel: _prep
	cd src/; make -f Makefile.packages wheel

pytest:
	cd src; pytest-3 -v -n auto --doctest-modules

update-golang-modules:
	cd src/; make -f Makefile.packages update-golang-modules

_prep:
	-mkdir src/dist;
	-ln -s src/build build;
	-ln -s src/dist dist;
	-ln -s src/deb_dist deb_dist;

clean:
	-rm build/*;
	-rm dist/*;
	-cd src; make clean
