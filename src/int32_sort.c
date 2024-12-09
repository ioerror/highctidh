#ifndef HIGHCTIDH_INT32_SORT_H
#define HIGHCTIDH_INT32_SORT_H

#include "config.h"

#if defined(ENABLE_ASM) && defined(__AVX2__)

// This is the original high-ctidh x86_64 sorting code
#define int32_MINMAX(a,b)			\
do { \
  int32 temp1; \
  __asm__( \
    ".att_syntax prefix\n\t" \
    "cmpl %1,%0\n\t" \
    "mov %0,%2\n\t" \
    "cmovg %1,%0\n\t" \
    "cmovg %2,%1\n\t" \
    : "+r"(a), "+r"(b), "=r"(temp1) \
    : \
    : "cc" \
  ); \
} while(0)

#else /* fallback to portable C code: */

// This is from the Public Domain release of djbsort-20190516
#include "int32_sort.h"
#define int32 int32_t

#define int32_MINMAX(a,b) do {			\
  register const int32_t big = (a > b ? a : b); \
  register const int32_t small = (a > b ? b : a); \
  a = small; \
  b = big; \
} while (0);

void int32_sort(int32 *x,long long n)
{
  long long top,p,q,r,i;

  if (n < 2) return;
  top = 1;
  while (top < n - top) top += top;

  for (p = top;p > 0;p >>= 1) {
    for (i = 0;i < n - p;++i)
      if (!(i & p))
        int32_MINMAX(x[i],x[i+p]);
    i = 0;
    for (q = top;q > p;q >>= 1) {
      for (;i < n - q;++i) {
        if (!(i & p)) {
          int32 a = x[i + p];
          for (r = q;r > p;r >>= 1)
            int32_MINMAX(a,x[i+r]);
      x[i + p] = a;
    }
      }
    }
  }
}

#endif /* end ENABLE_ASM */
#endif /* HIGHCTIDH_INT32_SORT_H */
