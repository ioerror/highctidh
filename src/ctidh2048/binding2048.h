#ifndef _BINDING_H
#define _BINDING_H

#include <stdlib.h>
#include <stdint.h>

void highctidh_2048_go_fillrandom(void *, void *, size_t);

__attribute__((weak))
void fillrandom_2048_custom(
  void *const outptr,
  const size_t outsz,
  const uintptr_t context)
{
  highctidh_2048_go_fillrandom((void *) context, outptr, outsz);
}
#endif
