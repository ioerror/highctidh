#include "highctidh.h"
#include <stdio.h>
#include <assert.h>
#include <string.h>

int
main(int argc, char **argv)
{
	(void) argc;
	(void) argv;

	/* p511 */
	highctidh_511_private_key sk_511_a;
	highctidh_511_public_key  pk_511_a;
	highctidh_511_csidh_private(&sk_511_a);
	assert(highctidh_511_csidh(&pk_511_a, &highctidh_511_base, &sk_511_a));

	highctidh_511_private_key sk_511_b;
	highctidh_511_public_key  pk_511_b;
	highctidh_511_csidh_private(&sk_511_b);
	assert(highctidh_511_csidh(&pk_511_b, &highctidh_511_base, &sk_511_b));
	highctidh_511_public_key shared_511_a;
	highctidh_511_public_key shared_511_b;
	assert(highctidh_511_csidh(&shared_511_a, &pk_511_a, &sk_511_b));
	assert(highctidh_511_csidh(&shared_511_b, &pk_511_b, &sk_511_a));
	assert(0 == memcmp(&shared_511_a, &shared_511_b, sizeof(shared_511_a)));

	/* p512 */
	highctidh_512_private_key sk_512_a;
	highctidh_512_public_key  pk_512_a;
	highctidh_512_csidh_private(&sk_512_a);
	assert(highctidh_512_csidh(&pk_512_a, &highctidh_512_base, &sk_512_a));

	highctidh_512_private_key sk_512_b;
	highctidh_512_public_key  pk_512_b;
	highctidh_512_csidh_private(&sk_512_b);
	assert(highctidh_512_csidh(&pk_512_b, &highctidh_512_base, &sk_512_b));
	highctidh_512_public_key shared_512_a;
	highctidh_512_public_key shared_512_b;
	assert(highctidh_512_csidh(&shared_512_a, &pk_512_a, &sk_512_b));
	assert(highctidh_512_csidh(&shared_512_b, &pk_512_b, &sk_512_a));
	assert(0 == memcmp(&shared_512_a, &shared_512_b, sizeof(shared_512_a)));

	/* p1024 */
	highctidh_1024_private_key sk_1024_a;
	highctidh_1024_public_key  pk_1024_a;
	highctidh_1024_csidh_private(&sk_1024_a);
	assert(highctidh_1024_csidh(&pk_1024_a, &highctidh_1024_base, &sk_1024_a));

	highctidh_1024_private_key sk_1024_b;
	highctidh_1024_public_key  pk_1024_b;
	highctidh_1024_csidh_private(&sk_1024_b);
	assert(highctidh_1024_csidh(&pk_1024_b, &highctidh_1024_base, &sk_1024_b));
	highctidh_1024_public_key shared_1024_a;
	highctidh_1024_public_key shared_1024_b;
	assert(highctidh_1024_csidh(&shared_1024_a, &pk_1024_a, &sk_1024_b));
	assert(highctidh_1024_csidh(&shared_1024_b, &pk_1024_b, &sk_1024_a));
	assert(0 == memcmp(&shared_1024_a, &shared_1024_b, sizeof(shared_1024_a)));

	/* p2048 */
	highctidh_2048_private_key sk_2048_a;
	highctidh_2048_public_key  pk_2048_a;
	highctidh_2048_csidh_private(&sk_2048_a);
	assert(highctidh_2048_csidh(&pk_2048_a, &highctidh_2048_base, &sk_2048_a));

	highctidh_2048_private_key sk_2048_b;
	highctidh_2048_public_key  pk_2048_b;
	highctidh_2048_csidh_private(&sk_2048_b);
	assert(highctidh_2048_csidh(&pk_2048_b, &highctidh_2048_base, &sk_2048_b));
	highctidh_2048_public_key shared_2048_a;
	highctidh_2048_public_key shared_2048_b;
	assert(highctidh_2048_csidh(&shared_2048_a, &pk_2048_a, &sk_2048_b));
	assert(highctidh_2048_csidh(&shared_2048_b, &pk_2048_b, &sk_2048_a));
	assert(0 == memcmp(&shared_2048_a, &shared_2048_b, sizeof(shared_2048_a)));

	return 0;
}
