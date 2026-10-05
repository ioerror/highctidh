#ifndef _BINDING_H
#define _BINDING_H

#ifdef CGONUTS

#include <stdlib.h>
#include <stdint.h>

#if 512 == BITS

void fillrandom_custom( void *const outptr, const size_t outsz, const uintptr_t context);
void highctidh_512_go_fillrandom(void *, void *, size_t);
#define NAMESPACEBITS(x) highctidh_512_##x
#define NAMESPACEGENERIC(x) highctidh_512_##x

__attribute__((weak))
void fillrandom_512_custom(
  void *const outptr,
  const size_t outsz,
  const uintptr_t context)
{
  highctidh_512_go_fillrandom((void *)context, outptr, outsz);
#if defined(__BYTE_ORDER__) && (__BYTE_ORDER__ == __ORDER_BIG_ENDIAN__)
  unsigned char *b = outptr;
  for (size_t i = 0; i + 4 <= outsz; i += 4) {
    unsigned char t = b[i];
    b[i] = b[i + 3];
    b[i + 3] = t;
    t = b[i + 1];
    b[i + 1] = b[i + 2];
    b[i + 2] = t;
  }
#endif
}
#endif

#endif

#endif
