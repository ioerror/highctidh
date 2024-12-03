//go:build arm64
#if HIGHCTIDH_PORTABLE == 0
.section .rodata

#include "uintbig_const_le_2048.s"

.section .text

#include "aarch64_macros.s"

.global highctidh_2048_uintbig_set
highctidh_2048_uintbig_set:
	SET	2048

.global highctidh_2048_uintbig_bit
highctidh_2048_uintbig_bit:
	BIT	2048

.global highctidh_2048_uintbig_add3
highctidh_2048_uintbig_add3:
	ADD3	2048

.global highctidh_2048_uintbig_sub3
highctidh_2048_uintbig_sub3:
	SUB3	2048

.global highctidh_2048_uintbig_mul3_64
highctidh_2048_uintbig_mul3_64:
	MUL3_64	2048
#endif
