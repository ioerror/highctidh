#include <string.h>
#include <assert.h>

#include "csidh.h"
#include "naidne.h"
#include "primes.h"
#include "random.h"

void csidh_private_withrng(private_key *priv, uintptr_t rng_context,
    ctidh_fillrandom rng_callback)
{
  memset(&priv->e, 0, sizeof(priv->e));
  long long pos = 0;
  long long w = 0;
  long long S = 0;
  for (long long b = 0;b < primes_batches;++b) {
    w = primes_batchsize[b];
    S = primes_batchbound[b];
    random_boundedl1(priv->e + pos,w,S, rng_context, rng_callback);
    pos += w;
  }
  assert(pos <= primes_num);
}

struct rng_bytes {
  ctidh_fillrandom *callback;
  uintptr_t context;
};

static void fillrandom_bytes(void *const outbuf, const size_t outsz,
    const uintptr_t context)
{
  const struct rng_bytes *r = (const struct rng_bytes *) context;
  unsigned char *b = outbuf;
  uint32_t w;
  r->callback(outbuf, outsz, r->context);
  for (size_t i = 0;i + 4 <= outsz;i += 4) {
    memcpy(&w, b + i, 4);
    w = le32toh(w);
    memcpy(b + i, &w, 4);
  }
}

void csidh_private_withrng_bytes(private_key *priv, uintptr_t rng_context,
    ctidh_fillrandom rng_callback)
{
  struct rng_bytes r = {rng_callback, rng_context};
  csidh_private_withrng(priv, (uintptr_t) &r, fillrandom_bytes);
}

void csidh_private(private_key *priv)
{
	csidh_private_withrng(priv, (uintptr_t) priv, ctidh_fillrandom_default);
}
