#if defined(__x86_64__)
#ifdef __clang__
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wlanguage-extension-token"
#endif
#define int32_MINMAX(a,b)			\
do { \
  int32 temp1; \
  asm( \
    "cmpl %1,%0\n\t" \
    "mov %0,%2\n\t" \
    "cmovg %1,%0\n\t" \
    "cmovg %2,%1\n\t" \
    : "+r"(a), "+r"(b), "=r"(temp1) \
    : \
    : "cc" \
  ); \
} while(0)
#else
#include <stdint.h>
#define int32_MINMAX(a,b) do { \
  const register int32_t big = (a > b ? a : b); \
  const register int32_t small = (a > b ? b : a); \
  a = small; \
  b = big; \
} while (0);
#ifdef __clang__
#pragma clang diagnostic pop
#endif
#endif
typedef int no_empty_translation_units; // -> "warning: ISO C forbids an empty translation unit"
