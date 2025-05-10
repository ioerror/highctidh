.section .rodata

#include "uintbig_const_le_512.s"

.section .text

#include "aarch64_macros.s"

.global highctidh_512_uintbig_set
highctidh_512_uintbig_set:
	SET	512

.global highctidh_512_uintbig_bit
highctidh_512_uintbig_bit:
	BIT	512

.global highctidh_512_uintbig_add3
highctidh_512_uintbig_add3:
	ADD3	512

.global highctidh_512_uintbig_sub3
highctidh_512_uintbig_sub3:
	SUB3	512

.global highctidh_512_uintbig_mul3_64
highctidh_512_uintbig_mul3_64:
	MUL3_64	512
