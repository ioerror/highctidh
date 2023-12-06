#include <stdlib.h>
#include <stdint.h>

__attribute__((weak))
void fillrandom_custom(
  void *const outptr,
  const size_t outsz,
  const uintptr_t context)
{
  extern void go_fillrandom(const uintptr_t context, void *const outptr, const size_t outsz);
}
