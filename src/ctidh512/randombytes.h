#ifndef randombytes_h
#define randombytes_h

#include <stdlib.h>
#include <stdint.h>
#if defined(CGONUTS)
#include "cgo.h"
#define randombytes NAMESPACEBITS(randombytes)
#endif


void randombytes(void *x, size_t l);

#endif
