AC_INIT([libhighctidh], 1.0.0, [authors], [libhighctidh])

# http://www.gnu.org/software/libtool/manual/html_node/Updating-version-info.html
# http://sourceware.org/autobook/autobook/autobook_91.html
lt_current="0"
lt_revision="0"
lt_age="0"

AM_INIT_AUTOMAKE([1.14 foreign dist-xz no-define subdir-objects])
LT_INIT([win32-dll])

AM_MAINTAINER_MODE([disable])
AC_CONFIG_SRCDIR([src/highctidh.h])
AC_CONFIG_MACRO_DIR([m4])
AM_CONFIG_HEADER([config.h])
m4_ifdef([AM_SILENT_RULES],[AM_SILENT_RULES([yes])])

AC_PREREQ([2.50])
AC_PROG_CC
AM_PROG_AS
AC_PROG_LIBTOOL
AC_C_INLINE
AC_DEFINE([_DEFAULT_SOURCE], 1, [htole64])

AC_MSG_CHECKING([for Windows])
case $host in
cygwin* | mingw* | pw32* | cegcc* | os2* | aix*)
	AC_MSG_RESULT([yes])
	LTLDFLAGS="-avoid-version -no-undefined"
	LTPROGLDFLAGS="-no-fast-install"
	;;
*)
	AC_MSG_RESULT([no])
	LTLDFLAGS="-version-info ${lt_current}:${lt_revision}:${lt_age}"
	LTPROGLDFLAGS="-no-install"
esac
AC_CHECK_HEADERS([sys/param.h])

AC_MSG_CHECKING([whether to use ubigint and fp assembly])
AC_ARG_ENABLE([asm],
	[AS_HELP_STRING([--enable-asm], [use ubigint and fp assembly for x86_64 BMI2 or armv8 SVE2])],
	[case "$enableval" in
		yes) asm=true ;;
		no) asm=false ;;
		*) AC_MSG_ERROR([bad value ${enableval} for --enable-asm]) ;;
	esac], [asm=false])
if test "x$asm" = "xtrue"; then
	AC_MSG_RESULT([yes])
	AC_DEFINE([ENABLE_ASM], 1, [optimized ubigint and fp assembly])
elif test "x$asm" = "xfalse"; then
	AC_MSG_RESULT([no])
fi
AM_CONDITIONAL(ENABLE_ASM, test "x$asm" = "xtrue")

AM_CFLAGS="-std=c11 -O2 -pedantic \
	-Wall -Wundef -Wunused -Wshadow -Wextra -Werror \
	-Wstrict-prototypes -Werror-implicit-function-declaration \
	-Wformat -Werror=format-security \
	-fwrapv -DGETRANDOM -fpie -fPIC -fPIE \
	-D_FORTIFY_SOURCE=2 -fstack-protector-strong \
	-D\"NAMESPACEGENERIC(x)=highctidh_\\#\\#x\" \
	-D\"NAMESPACE2(SIZE,x)=highctidh_\\#\\#SIZE\\#\\#_\\#\\#x\" \
	-D\"NAMESPACE1(SIZE,x)=NAMESPACE2(SIZE,x)\" \
	-D\"NAMESPACEBITS(x)=NAMESPACE1(BITS,x)\""

AC_SUBST([AM_CFLAGS])
AC_SUBST([LTLDFLAGS])
AC_SUBST([LTPROGLDFLAGS])

AC_CONFIG_FILES([libhighctidh_511.pc])
AC_CONFIG_FILES([libhighctidh_512.pc])
AC_CONFIG_FILES([libhighctidh_1024.pc])
AC_CONFIG_FILES([libhighctidh_2048.pc])
AC_CONFIG_FILES([libhighctidh.pc])
AC_CONFIG_FILES([Makefile])
AC_OUTPUT
