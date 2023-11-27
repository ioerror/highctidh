#ifndef HIGHCTIDH_H
#define HIGHCTIDH_H

/*
 * This header file describes the public API of libhighctidh for each
 * field size.
 *
 * See csidh.h for descriptions of functions and typedefs.
 *
 */

#include <stdint.h>
#include <stddef.h>

typedef void ((ctidh_fillrandom)(
  void *const outbuf,
  const size_t outsz,
  const uintptr_t context));

ctidh_fillrandom ctidh_fillrandom_default;

/* 511 */
typedef struct highctidh_511_uintbig {
    uint64_t c[((511 +63)/64)];
} highctidh_511_uintbig;

typedef struct highctidh_511_fp {
    highctidh_511_uintbig x;
} highctidh_511_fp;

typedef struct highctidh_511_private_key {
    int8_t e[74];
} highctidh_511_private_key;

typedef struct highctidh_511_public_key {
    highctidh_511_fp A;
} highctidh_511_public_key;

extern const highctidh_511_public_key highctidh_511_base;
void highctidh_511_public_key_from_bytes(highctidh_511_public_key *const pk, const char *const input);
void highctidh_511_public_key_to_bytes(char *const output, const highctidh_511_public_key *const pk);

void highctidh_511_csidh_private_withrng(highctidh_511_private_key *priv, uintptr_t rng_context, ctidh_fillrandom rng_callback);
void highctidh_511_csidh_private(highctidh_511_private_key *const priv);
_Bool highctidh_511_csidh(highctidh_511_public_key *out, highctidh_511_public_key const *in, highctidh_511_private_key const *priv);
_Bool highctidh_511_validate(highctidh_511_public_key const *in);
void highctidh_511_action(highctidh_511_public_key *out, highctidh_511_public_key const *in, highctidh_511_private_key const *priv);


/* 512 */

typedef struct highctidh_512_uintbig {
    uint64_t c[((512 +63)/64)];
} highctidh_512_uintbig;
typedef struct highctidh_512_fp {
    highctidh_512_uintbig x;
} highctidh_512_fp;
typedef struct highctidh_512_private_key {
    int8_t e[74];
} highctidh_512_private_key;
typedef struct highctidh_512_public_key {
    highctidh_512_fp A;
} highctidh_512_public_key;
extern const highctidh_512_public_key highctidh_512_base;
void highctidh_512_public_key_from_bytes(highctidh_512_public_key *const pk, const char *const input);
void highctidh_512_public_key_to_bytes(char *const output, const highctidh_512_public_key *const pk);
void highctidh_512_csidh_private_withrng(highctidh_512_private_key *priv, uintptr_t rng_context, ctidh_fillrandom rng_callback);
void highctidh_512_csidh_private(highctidh_512_private_key *const priv);
_Bool highctidh_512_csidh(highctidh_512_public_key *out, highctidh_512_public_key const *in, highctidh_512_private_key const *priv);
_Bool highctidh_512_validate(highctidh_512_public_key const *in);
void highctidh_512_action(highctidh_512_public_key *out, highctidh_512_public_key const *in, highctidh_512_private_key const *priv);


/* 1024 */

typedef struct highctidh_1024_uintbig {
    uint64_t c[((1024 +63)/64)];
} highctidh_1024_uintbig;
typedef struct highctidh_1024_fp {
    highctidh_1024_uintbig x;
} highctidh_1024_fp;
typedef struct highctidh_1024_private_key {
    int8_t e[130];
} highctidh_1024_private_key;
typedef struct highctidh_1024_public_key {
    highctidh_1024_fp A;
} highctidh_1024_public_key;
extern const highctidh_1024_public_key highctidh_1024_base;
void highctidh_1024_public_key_from_bytes(highctidh_1024_public_key *const pk, const char *const input);
void highctidh_1024_public_key_to_bytes(char *const output, const highctidh_1024_public_key *const pk);
typedef void ((ctidh_fillrandom)(
  void *const outbuf,
  const size_t outsz,
  const uintptr_t context));

void highctidh_1024_csidh_private_withrng(highctidh_1024_private_key *priv, uintptr_t rng_context, ctidh_fillrandom rng_callback);
void highctidh_1024_csidh_private(highctidh_1024_private_key *const priv);
_Bool highctidh_1024_csidh(highctidh_1024_public_key *out, highctidh_1024_public_key const *in, highctidh_1024_private_key const *priv);
_Bool highctidh_1024_validate(highctidh_1024_public_key const *in);
void highctidh_1024_action(highctidh_1024_public_key *out, highctidh_1024_public_key const *in, highctidh_1024_private_key const *priv);


/* 2048 */

typedef struct highctidh_2048_uintbig {
    uint64_t c[((2048 +63)/64)];
} highctidh_2048_uintbig;
typedef struct highctidh_2048_fp {
    highctidh_2048_uintbig x;
} highctidh_2048_fp;
typedef struct highctidh_2048_private_key {
    int8_t e[231];
} highctidh_2048_private_key;
typedef struct highctidh_2048_public_key {
    highctidh_2048_fp A;
} highctidh_2048_public_key;
extern const highctidh_2048_public_key highctidh_2048_base;
void highctidh_2048_public_key_from_bytes(highctidh_2048_public_key *const pk, const char *const input);
void highctidh_2048_public_key_to_bytes(char *const output, const highctidh_2048_public_key *const pk);
void highctidh_2048_csidh_private_withrng(highctidh_2048_private_key *priv, uintptr_t rng_context, ctidh_fillrandom rng_callback);
void highctidh_2048_csidh_private(highctidh_2048_private_key *const priv);
_Bool highctidh_2048_csidh(highctidh_2048_public_key *out, highctidh_2048_public_key const *in, highctidh_2048_private_key const *priv);
_Bool highctidh_2048_validate(highctidh_2048_public_key const *in);
void highctidh_2048_action(highctidh_2048_public_key *out, highctidh_2048_public_key const *in, highctidh_2048_private_key const *priv);

#endif /* HIGHCTIDH_H */
