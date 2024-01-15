#ifdef CGONUTS
#include <stdlib.h>
#include <stdint.h>
#include "binding512.h"

__attribute__((weak))
void fillrandom_custom(
  void *const outptr,
  const size_t outsz,
  const uintptr_t context)
{
  highctidh_512_go_fillrandom((void *)context, outptr, outsz);
}
#endif
