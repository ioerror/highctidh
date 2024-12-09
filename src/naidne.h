#ifndef HIGHCTIDH_NAIDNE_H
#define HIGHCTIDH_NAIDNE_H

#include <config.h>

#if defined(HAVE_ENDIAN_H)
#include <endian.h>
#elif defined(HAVE_SYS_ENDIAN_H)
#include <sys/endian.h>
#elif defined(HAVE_MACHINE_ENDIAN_H) && defined(HAVE_LIBKERN_OSBYTEORDER_H)
#include <sys/types.h>
#include <machine/endian.h>
#include <libkern/OSByteOrder.h>
#define htole32(x) OSSwapHostToLittleInt32(x)
#define htole64(x) OSSwapHostToLittleInt64(x)
#define le32toh(x) OSSwapLittleToHostInt32(x)
#define le64toh(x) OSSwapLittleToHostInt64(x)
#elif defined(HAVE_SYS_BYTEORDER_H)
#include <sys/byteorder.h>
#define htole32(x) LE_32(x)
#define htole64(x) LE_64(x)
#define le32toh(x) LE_32(x)
#define le64toh(x) LE_64(x)
#elif (defined(__Windows__) || defined(_WIN64) || defined(_WIN32))
#include <stdlib.h>
#ifdef HAVE_SYS_PARAM_H
#include <sys/param.h>
#endif
#if (BYTE_ORDER == LITTLE_ENDIAN)
#define htole32(x) (x)
#define le32toh(x) (x)
#define htole64(x) (x)
#define le64toh(x) (x)
#elif (BYTE_ORDER == BIG_ENDIAN)
#define htole32(x) _byteswap_ulong(x)
#define le32toh(x) _byteswap_ulong(x)
#define htole64(x) _byteswap_uint64(x)
#define le64toh(x) _byteswap_uint64(x)
#endif
#else
#error No known endian function support.
#endif

#endif
