#ifndef _BINDING_H
#define _BINDING_H

#ifdef CGONUTS

#include <stdlib.h>
#include <stdint.h>

#if 512 == BITS

void fillrandom_512_custom( void *const outptr, const size_t outsz, const uintptr_t context);
void highctidh_512_go_fillrandom(void *, void *, size_t);
#define NAMESPACEBITS(x) highctidh_512_##x
#define NAMESPACEGENERIC(x) highctidh_512_##x

#endif

#endif

#endif
