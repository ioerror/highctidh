#ifndef _BINDING_H
#define _BINDING_H

#ifdef CGONUTS

#include <stdlib.h>
#include <stdint.h>

#if 511 == BITS

void fillrandom_custom( void *const outptr, const size_t outsz, const uintptr_t context);
void go_fillrandom(void *, void *, size_t);
#define NAMESPACEBITS(x) highctidh_511_##x
#define NAMESPACEGENERIC(x) highctidh_##x

#endif

#endif

#endif