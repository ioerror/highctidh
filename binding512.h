#ifndef _BINDING_H
#define _BINDING_H

#ifdef CGONUTS

#if 512 == BITS

#define NAMESPACEBITS(x) highctidh_512_##x
#define NAMESPACEGENERIC(x) highctidh_##x

#endif

#endif

#endif
