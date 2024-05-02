#ifdef CGONUTS
#include <stdlib.h>
#include <stdint.h>
#include "binding511.h"

#define go_fillrandom NAMESPACEBITS(go_fillrandom)

__attribute__((weak))
void fillrandom_511_custom(
  void *const outptr,
  const size_t outsz,
  const uintptr_t context)
{
  go_fillrandom((void *)context, outptr, outsz);
}

#endif
