#ifndef randombytes_h
#define randombytes_h

#include <stdlib.h>
#include <stdint.h>
#if defined(CGONUTS)
#include "cgo.h"
#define randombytes NAMESPACEBITS(randombytes)
#endif


void randombytes(void *x, size_t l);
#if (defined(__Windows__) || defined(__WIN64))
ssize_t getrandom(void buf, size_t buflen);
#endif

#endif
