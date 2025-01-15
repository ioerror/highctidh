/*
 * This program shows basic usage with the highctidh csidh API.
 *
 * Compile this example program for each key size using the system installed
 * headers and linked against the installed shared library:
 *
 *  gcc -DBITS=511 -D'NAMESPACEBITS(x)=highctidh_511_##x' \
 *  -D'NAMESPACEGENERIC(x)=highctidh_##x' -o example-ctidh511 example-ctidh.c \
 *  -l highctidh_511
 *
 *  gcc -DBITS=512 -D'NAMESPACEBITS(x)=highctidh_512_##x' \
 *  -D'NAMESPACEGENERIC(x)=highctidh_##x' -o example-ctidh512 example-ctidh.c \
 *  -l highctidh_512
 *
 *  gcc -DBITS=1024 -D'NAMESPACEBITS(x)=highctidh_1024_##x' \
 *  -D'NAMESPACEGENERIC(x)=highctidh_##x' -o example-ctidh1024 example-ctidh.c \
 *  -l highctidh_1024
 *
 *  gcc -DBITS=2048 -D'NAMESPACEBITS(x)=highctidh_2048_##x' \
 *  -D'NAMESPACEGENERIC(x)=highctidh_##x' -o example-ctidh2048 example-ctidh.c \
 *  -l highctidh_2048
 *
 * */

#include <stdio.h>
#include <libhighctidh/csidh.h>

void print_hex_key(void *k, unsigned int l)
{ 
  printf("0x");
  for (unsigned int i = 0; i < l; i++)
  {
    printf("%02x", i[(unsigned char *) k]);
  }
  printf("\n");
}

int main(void);
int main(void)
{
  private_key sk_a, sk_b;
  public_key pk_a, pk_b;
  public_key s_a, s_b;
  bool ok = 0;

  setvbuf(stdout, NULL, _IOLBF, 4096);

  printf("CTIDH %i vector example\n\n", BITS);

  printf("Generating Alice's private_key (%zu bytes):\n", sizeof sk_a);
  csidh_private(&sk_a);
  print_hex_key(&sk_a, sizeof sk_a);
  printf("Generating Alice's public_key (%zu bytes):\n", sizeof pk_a);
  ok = csidh(&pk_a, &base, &sk_a);
  if (!validate(&pk_a))
  {
    printf("Invalid public key:\n");
  }
  print_hex_key(&pk_a, sizeof pk_a);
  printf("Result: %i\n", ok);
  printf("\n");

  printf("Generating Bob's private_key (%zu bytes):\n", sizeof sk_b);
  csidh_private(&sk_b);
  print_hex_key(&sk_b, sizeof sk_b);
  printf("Generating Bob's public_key (%zu bytes):\n", sizeof pk_b);
  ok = csidh(&pk_b, &base, &sk_b);
  if (!validate(&pk_a))
  {
    printf("Invalid public key:\n");
  }
  print_hex_key(&pk_b, sizeof pk_b);
  printf("Result: %i\n", ok);
  printf("\n");

  printf("Computing DH for Alice (%zu bytes):\n", sizeof s_a);
  ok = csidh(&s_a, &pk_b, &sk_a);
  print_hex_key(&s_a, sizeof s_a);
  printf("Result: %i\n", ok);
  printf("\n");

  printf("Computing DH for Bob (%zu bytes):\n", sizeof s_b);
  ok = csidh(&s_b, &pk_a, &sk_b);
  print_hex_key(&s_b, sizeof s_b);
  printf("Result: %i\n", ok);
  printf("\n");

  printf("Shared keys ...");
  if (!memcmp(&s_a, &s_b, sizeof s_a))
  {
    printf(" match\n");
    ok = 0;
  } else {
    printf(" do not match\n");
    ok = 1;
  }
  return ok;
}
