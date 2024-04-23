#ifndef HIGHCTIDH_NAIDNE_H
#define HIGHCTIDH_NAIDNE_H

#if defined(__linux__)
#include <endian.h>
#elif defined(__FreeBSD__)
#include <sys/types.h>
#include <sys/endian.h>
#elif defined(__APPLE__)
#include <sys/types.h>
#include <machine/endian.h>
#include <libkern/OSByteOrder.h>
#define htole32(x) OSSwapHostToLittleInt32(x)
#define htole64(x) OSSwapHostToLittleInt64(x)
#define le32toh(x) OSSwapLittleToHostInt32(x)
#define le64toh(x) OSSwapLittleToHostInt64(x)
#elif defined(__sun)
#include <sys/byteorder.h>
#define htole32(x) LE_32(x)
#define htole64(x) LE_64(x)
#define le32toh(x) LE_32(x)
#define le64toh(x) LE_64(x)
#else
#include <endian.h>
#endif

#endif
