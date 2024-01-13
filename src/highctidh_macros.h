#ifndef CSIDH_ALL_H
#define CSIDH_ALL_H

/*
 * This file includes csidh.h for each key size, so it can be used
 * stand-alone without external macro definitions.
 *
 * It is used to generate highctidh.h, with an automatic pre-processing step:
 * gcc -I. -DGETRANDOM -DHIGHCTIDH_PORTABLE -E highctidh_macros.h | \
 * awk '/^#/ { if($3 ~ /"[^/]/){nice=1} else {nice = 0}} nice && /^[^#]/ {print}'
 * and some manual cleanup because the oneliner above isn't perfect.
 *
 */

#ifdef BITS
#error "don't include highctidh.h with the BITS macro defined"
#endif

#define NAMESPACEGENERIC(x) highctidh_##x

#define BITS 511
#include "csidh_all_clearnamespaces.h"
#include "csidh.h"

#undef BITS
#define BITS 512
#include "csidh_all_clearnamespaces.h"
#include "csidh.h"

#undef BITS
#define BITS 1024
#include "csidh_all_clearnamespaces.h"
#include "csidh.h"

#undef BITS
#define BITS 2048
#include "csidh_all_clearnamespaces.h"
#include "csidh.h"
#undef BITS

#endif /* CSIDH_ALL_H */
