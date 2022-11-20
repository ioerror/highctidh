#ifdef BITS
#error "don't include highctidh.h with the BITS macro defined"
#endif

#ifndef CSIDH_ALL_H
#define CSIDH_ALL_H

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
