library:
	cd src; make;
	cp -v src/*.so dist/;

packages:
	cd src/; make -f Makefile.packages packages;

deb:
	cd src/; make -f Makefile.packages deb

wheel:
	cd src/; make -f Makefile.packages wheel

update-golang-modules:
	cd src/; make -f Makefile.packages update-golang-modules

clean:
	-rm build/*;
	-rm dist/*;
	-cd src; make clean
