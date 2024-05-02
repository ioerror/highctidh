#ifndef _BINDING_H
#define _BINDING_H

#ifdef CGONUTS

#include <stdlib.h>
#include <stdint.h>

#if 511 == BITS

#if (!defined(__Windows__) || !defined(__WIN64))
void fillrandom_511_custom( void *const outptr, const size_t outsz, const uintptr_t context);
void highctidh_511_go_fillrandom(void *, void *, size_t);
#endif
#define NAMESPACEBITS(x) highctidh_511_##x
#define NAMESPACEGENERIC(x) highctidh_511_##x

#endif

#endif

#endif
