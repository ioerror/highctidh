#define NAMESPACEGENERIC(x) highctidh_ ## x
#define NAMESPACE2(SIZE,x) highctidh_ ## SIZE ## _ ## x
#define NAMESPACE1(SIZE,x) NAMESPACE2(SIZE,x)
#define NAMESPACEBITS(x) NAMESPACE1(BITS,x)
