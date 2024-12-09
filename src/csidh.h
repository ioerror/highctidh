#ifndef CSIDH_H
#define CSIDH_H

#include "naidne.h"
#include "uintbig.h"
#include "fp.h"
#include "mont.h"
#include "primes.h"
#include "csidh_namespace.h"
#include "random.h"

extern long long csidh_stattried[primes_batches];
extern long long csidh_statsucceeded[primes_batches];

typedef struct private_key {
    int8_t e[primes_num];
} private_key;

typedef struct public_key {
    fp A; /* Montgomery coefficient: represents y^2 = x^3 + Ax^2 + x */
} public_key;

extern const public_key base;

/*
 * Initialize a public_key from a byte array of length sizeof(public_key).
 * This is required to ensure interoperability between
 * little- and big-endian systems, since the limbs internally
 * must be in host/native order.
 */
void public_key_from_bytes(public_key *const pk, const char *const input);

/*
 * Serialize a public_key to a byte array of length sizeof(public_key).
 * This is required to ensure interoperability between
 * little- and big-endian systems, since the limbs internally
 * must be in host/native order.
 */
void public_key_to_bytes(char *const output, const public_key *const pk);

/*
 * generate a new private key using rng_callback and write the result to (priv).
 * (context) is passed as context to the (rng_callback).
 */
void csidh_private_withrng(private_key *priv, uintptr_t rng_context, ctidh_fillrandom rng_callback);

/*
 * Generate a new private key and write the result to (priv).
 */
void csidh_private(private_key *const priv);

/*
 * Evaluates the group action (the "Diffie-Hellman"-like function).
 * Returns:
 * false: when (in) is not a valid public key. (out) filled with random bytes.
 * true: when (in) is a valid key. (out) is the resulting field element.
 */
bool csidh(public_key *out, public_key const *in, private_key const *priv);

int validate_cutofforder_v2(uintbig *order,const fp *P,const fp *A);

/*
 * Validates a public_key and returns true when valid; false when invalid.
 */
bool validate(public_key const *in);

/*
 * Evaluates the group action WITHOUT validating the (in) public_key.
 * This function can be used instead of csidh() when the public_key has already
 * been validated.
 */
void action(public_key *out, public_key const *in, private_key const *priv);

#endif
