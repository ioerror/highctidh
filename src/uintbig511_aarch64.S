.section .rodata

#include "uintbig_const_le_511.s"

.section .text

#include "aarch64_macros.s"

.global highctidh_511_uintbig_set
highctidh_511_uintbig_set:
	SET	512

.global highctidh_511_uintbig_bit
.hidden bit512
highctidh_511_uintbig_bit:
	cmp	x1, #511-1
	b.ls	bit512
	mov	x0, #0
	ret
bit512:
	BIT	512

.global highctidh_511_uintbig_add3
highctidh_511_uintbig_add3:
	ADD3	512

.global highctidh_511_uintbig_sub3
highctidh_511_uintbig_sub3:
	SUB3	512

.global highctidh_511_uintbig_mul3_64
highctidh_511_uintbig_mul3_64:
	MUL3_64	512
