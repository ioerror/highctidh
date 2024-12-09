#ifndef RANDOM_H
#define RANDOM_H

#include <inttypes.h>
#include <stdint.h>
#include <stddef.h>

/*
 * The (ctidh_fillrandom) function signature for custom rng implementations.
 * The (context) parameter can be used to implement thread-safe deterministic
 * CSPRNGs, when (context) is unique for parallel calls.
 *
 * Note that to achieve reproducible public_key derivation, the rng must write
 * the random bytes as an array of int32_t values with host-order/native
 * endianness. ie when it writes the following on a little-endian machine:
 * AA BB CC DD EE FF GG HH 11 22 33 44
 * it must write this on a big-endian machine:
 * DD CC BB AA HH GG FF EE 44 33 22 11
 * This means care must be taken to byteswap when using e.g. HKDF (whose
 * output state is usually standardized to be written in little-endian).
 */
typedef void ((ctidh_fillrandom)(
  void *const outbuf,
  const size_t outsz,
  const uintptr_t context));

/*
 * The default RNG calls getrandom() or reads from /dev/urandom
 */
ctidh_fillrandom ctidh_fillrandom_default;

// set up e[0]..e[w-1] having l1 norm at most S, assuming S<128, w<128
void random_boundedl1(int8_t *e,long long w,long long S, uintptr_t rng_context, ctidh_fillrandom rng_callback);

// return -1 with probability num/den, 0 with probability 1-num/den
// assuming 0 <= num <= den, 0 < den < 2^63
int64_t random_coin(uint64_t num,uint64_t den);

void random_boundedl1(int8_t *e, const long long w,const long long S, uintptr_t rng_context, ctidh_fillrandom rng_callback);

#endif
