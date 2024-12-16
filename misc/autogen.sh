#!/bin/sh
aclocal --force -I m4
libtoolize --copy --force
aclocal --force -I m4
autoconf --force
autoheader --force
automake --add-missing --copy --force-missing
