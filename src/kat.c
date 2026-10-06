#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "csidh.h"

struct splitmix64 {
	uint64_t state;
	unsigned char buf[8];
	size_t have;
};

static void
splitmix64_fill(void *const out, const size_t outsz, const uintptr_t context)
{
	struct splitmix64 *r = (struct splitmix64 *)context;
	unsigned char *p = out;
	for (size_t i = 0; i < outsz; i++) {
		if (r->have == 0) {
			uint64_t z = (r->state += 0x9e3779b97f4a7c15ULL);
			z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9ULL;
			z = (z ^ (z >> 27)) * 0x94d049bb133111ebULL;
			z ^= z >> 31;
			for (size_t j = 0; j < 8; j++)
				r->buf[j] = (unsigned char)(z >> (8 * j));
			r->have = 8;
		}
		p[i] = r->buf[8 - r->have--];
	}
	for (size_t i = 0; i + 4 <= outsz; i += 4) {
		uint32_t w;
		memcpy(&w, p + i, 4);
		w = le32toh(w);
		memcpy(p + i, &w, 4);
	}
}

static void
print_hex(const char *name, const void *buf, size_t len)
{
	const unsigned char *p = buf;
	printf("%s = ", name);
	for (size_t i = 0; i < len; i++)
		printf("%02x", p[i]);
	printf("\n");
}

static void
print_pk(const char *name, const public_key *pk)
{
	char b[sizeof(pk->A.x.c)];
	public_key_to_bytes(b, pk);
	print_hex(name, b, sizeof(b));
}

static void
keygen(private_key *sk, uint64_t seed)
{
	struct splitmix64 r = {seed, {0}, 0};
	csidh_private_withrng(sk, (uintptr_t)&r, splitmix64_fill);
}

static int
wire_action(public_key *out, const public_key *in, const private_key *sk)
{
	public_key wire;
	char b[sizeof(in->A.x.c)];
	public_key_to_bytes(b, in);
	public_key_from_bytes(&wire, b);
	return csidh(out, &wire, sk);
}

int
main(int argc, char **argv)
{
	private_key sk_a, sk_b, sk_f;
	public_key pk_a, pk_b, ss, bl;
	uint64_t seed[3];

	if (argc != 4)
		return 2;
	for (int i = 0; i < 3; i++)
		seed[i] = strtoull(argv[i + 1], NULL, 10);
	keygen(&sk_a, seed[0]);
	keygen(&sk_b, seed[1]);
	keygen(&sk_f, seed[2]);
	if (!wire_action(&pk_a, &base, &sk_a) || !wire_action(&pk_b, &base, &sk_b) ||
	    !wire_action(&ss, &pk_b, &sk_a) || !wire_action(&bl, &pk_a, &sk_f))
		return 1;
	printf("size = %d\nrng = splitmix64\n", BITS);
	printf("seed_a = %" PRIu64 "\n", seed[0]);
	print_hex("sk_a", sk_a.e, sizeof(sk_a.e));
	print_pk("pk_a", &pk_a);
	printf("seed_b = %" PRIu64 "\n", seed[1]);
	print_hex("sk_b", sk_b.e, sizeof(sk_b.e));
	print_pk("pk_b", &pk_b);
	printf("seed_f = %" PRIu64 "\n", seed[2]);
	print_hex("sk_f", sk_f.e, sizeof(sk_f.e));
	print_pk("ss_a_b", &ss);
	print_pk("blind_f_a", &bl);
	return 0;
}
