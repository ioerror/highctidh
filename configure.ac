AC_INIT([libhighctidh], [1.0.0], [authors], [libhighctidh])
AM_INIT_AUTOMAKE([1.14 foreign no-define no-dependencies subdir-objects])
AM_MAINTAINER_MODE([disable])
AC_CONFIG_SRCDIR([src/highctidh.h])
AC_CONFIG_MACRO_DIR([m4])
AC_CONFIG_HEADERS([config.h])
m4_ifdef([AM_SILENT_RULES], AM_SILENT_RULES([yes]))
AC_PREREQ([2.69])
AC_PROG_CC
AM_PROG_AS

AC_C_INLINE
AC_C_RESTRICT

AC_DEFINE([_DEFAULT_SOURCE], [1], [le64toh])

AM_CFLAGS="-std=c11 -O2 -fwrapv -fno-lto"

AC_CHECK_TYPE([__int128], [], [
	AC_MSG_CHECKING([whether $CC supports -fforce-enable-int128])
	saved_CFLAGS="$CFLAGS"
	CFLAGS="$AM_CFLAGS -fforce-enable-int128"
	AC_COMPILE_IFELSE([AC_LANG_SOURCE([int main(void) { __int128 i; (void)i; return 0; }])], [
		AC_MSG_RESULT([yes])
		AM_CFLAGS="$AM_CFLAGS -fforce-enable-int128"
	], [
		AC_MSG_RESULT([no])
		AC_MSG_ERROR([can not build without __int128 support])
	])
	CFLAGS="$saved_CFLAGS"
])

# https://www.gnu.org/software/libtool/manual/html_node/Updating-version-info.html
# https://sourceware.org/autobook/autobook/autobook_91.html
lt_cur="0"
lt_rev="0"
lt_age="0"

AS_IF([test "x$CC" != xemcc], [
	AC_MSG_CHECKING([for Windows])
	case $host in
	cygwin* | *-mingw* | pw32* | cegcc* | os2* | aix*)
		AC_MSG_RESULT([yes])
		CCASFLAGS="-DWINDOWS=1"
		LTLDFLAGS="-avoid-version -no-undefined"
		LTPROGLDFLAGS="-no-fast-install"
		;;
	*)
		AC_MSG_RESULT([no])
		CCASFLAGS="-DWINDOWS=0"
		AM_CFLAGS="$AM_CFLAGS -fstack-protector-strong"
		LTLDFLAGS="-version-info ${lt_cur}:${lt_rev}:${lt_age}"
		LTPROGLDFLAGS="-no-install"
		;;
	esac

	AC_ENABLE_SHARED
	AC_ENABLE_STATIC
], [
	AC_DISABLE_SHARED
	AC_ENABLE_STATIC
])
LT_INIT([disable-fast-install])

AC_CHECK_HEADERS([endian.h fcntl.h libkern/OSByteOrder.h machine/endian.h sys/byteorder.h sys/endian.h sys/param.h])

AC_ARG_WITH([fortify],
	AS_HELP_STRING([--with-fortify], [_FORTIFY_SOURCE value @<:@default=2@:>@]), [
		AS_IF([test "x$withval" = "xyes"], [withval=2])
		case "$withval" in
		0|no) ;;
		1|2|3) AM_CFLAGS="$AM_CFLAGS -D_FORTIFY_SOURCE=$withval" ;;
		*) AC_MSG_ERROR([invalid value '$withval' for --with-fortify]) ;;
		esac
	])

AC_MSG_CHECKING([whether to use ubigint and fp assembly])
AC_ARG_ENABLE([asm],
	AS_HELP_STRING([--enable-asm], [use BMI2/SVE2 ubigint and fp assembly]), [
		case "$enableval" in
		yes|no) ;;
		*) AC_MSG_ERROR([invalid value '$enableval' for --enable-asm]) ;;
		esac
	],
	[enable_asm=no]
)
AC_MSG_RESULT([$enable_asm])
AS_IF([test "x$enable_asm" = xyes], [
	AS_IF([test "x${CC%%emcc*}" != "x$CC"], [
		AC_MSG_ERROR([can not use assembly with emcc]),
		AC_DEFINE([ENABLE_ASM], [1], [optimized ubigint and fp assembly])
	])
])
AM_CONDITIONAL([ENABLE_ASM], [test "x$enable_asm" = xyes])

AC_ARG_WITH([getrandom],
	AS_HELP_STRING([--with-getrandom], [use getrandom() @<:@default=check@:>@]))
AS_IF([test "x$with_getrandom" != xno], [
	AC_MSG_CHECKING([for getrandom])
	AC_LINK_IFELSE([AC_LANG_SOURCE([
		#include <sys/random.h>
		int main() { return getrandom((void *)0, 0U, 0U); }
	])], [
		AC_DEFINE([HAVE_GETRANDOM], [1], [Define to 1 if you have the `getrandom' function.])
		AC_MSG_RESULT([yes])
	], [
		AC_MSG_RESULT([no])
		AS_IF([test "x$with_getrandom" = xyes], [
			AC_MSG_ERROR([forced --with-getrandom, but no getrandom() detected])
		])
	])
])

AC_FUNC_MALLOC
AC_CHECK_FUNCS([memset sqrt])

AC_ARG_ENABLE([assert],
	AS_HELP_STRING([--disable-assert], [make assert() do nothing]))
AS_IF([test "x$enable_assert" != "xyes"], [
	AC_DEFINE([NDEBUG], [1], [assert() generates no code])
])

AC_ARG_ENABLE([debug],
	AS_HELP_STRING([--enable-debug], [Add -fno-omit-frame-pointer -g3 to CFLAGS]))
AS_IF([test "x$enable_debug" != "xno"], [AM_CFLAGS="$AM_CFLAGS -fno-omit-frame-pointer"])

AM_CFLAGS="$AM_CFLAGS -pedantic -Wall -Wundef -Wunused -Wshadow -Wextra -Werror \
 -Wstrict-prototypes -Werror-implicit-function-declaration \
 -Wformat -Werror=format-security $CFLAGS"

AS_IF([test "x$enable_debug" != "xno"], [AM_CFLAGS="$AM_CFLAGS -g3"])

AC_SUBST([AM_CFLAGS])
AC_SUBST([CCASFLAGS])
AC_SUBST([LTLDFLAGS])
AC_SUBST([LTPROGLDFLAGS])

AC_CONFIG_FILES([libhighctidh_511.pc])
AC_CONFIG_FILES([libhighctidh_512.pc])
AC_CONFIG_FILES([libhighctidh_1024.pc])
AC_CONFIG_FILES([libhighctidh_2048.pc])
AC_CONFIG_FILES([libhighctidh.pc])
AC_CONFIG_FILES([Makefile])
AC_OUTPUT
