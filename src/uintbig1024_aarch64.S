.section .rodata

#include "uintbig_const_le_1024.s"

.section .text

#include "aarch64_macros.s"

.global highctidh_1024_uintbig_set
highctidh_1024_uintbig_set:
	SET	1024

.global highctidh_1024_uintbig_bit
highctidh_1024_uintbig_bit:
	BIT	1024

.global highctidh_1024_uintbig_add3
highctidh_1024_uintbig_add3:
	ADD3	1024

.global highctidh_1024_uintbig_sub3
highctidh_1024_uintbig_sub3:
	SUB3	1024

.global highctidh_1024_uintbig_mul3_64
highctidh_1024_uintbig_mul3_64:
	MUL3_64	1024
