#include "aarch64_macros.s"

.set plimbs, 16
.set pbytes, plimbs*8

.section .rodata

#include "fp_const_le_1024.s"

.section .text
.p2align 4,,15

.extern highctidh_1024_uintbig_p
.hidden highctidh_1024_uintbig_p


.global highctidh_1024_fp_copy
highctidh_1024_fp_copy:
	LD64x8	x1, 0, x2, x3, x4, x5, x6, x7, x8, x9
	LD64x8	x1, 64, x10, x11, x12, x13, x14, x15, x16, x17
	ST64x8	x0, 0, x2, x3, x4, x5, x6, x7, x8, x9
	ST64x8	x0, 64, x10, x11, x12, x13, x14, x15, x16, x17
	ret


.global highctidh_1024_fp_cmov
highctidh_1024_fp_cmov:
	and	x2, x2, #1
	neg	x2, x2

	.set k, 0
	.rept plimbs/2
		ldp	x3, x4, [x0, #8*k]
		ldp	x5, x6, [x1, #8*k]

		eor	x5, x5, x3
		and	x5, x5, x2
		eor	x3, x3, x5

		eor	x6, x6, x4
		and	x6, x6, x2
		eor	x4, x4, x6

		stp	x3, x4, [x0, #8*k]

		.set k, k+2
	.endr
	ret


.global highctidh_1024_fp_cswap
highctidh_1024_fp_cswap:
	and	x2, x2, #1
	neg	x2, x2

	.set k, 0
	.rept plimbs/2
		ldp	x3, x4, [x0, #8*k]
		ldp	x5, x6, [x1, #8*k]

		eor	x7, x3, x5
		and	x7, x7, x2
		eor	x3, x3, x7
		eor	x5, x5, x7

		eor	x8, x4, x6
		and	x8, x8, x2
		eor	x4, x4, x8
		eor	x6, x6, x8

		stp	x3, x4, [x0, #8*k]
		stp	x5, x6, [x1, #8*k]

		.set k, k+2
	.endr
	ret


.hidden reduce_once
reduce_once:
	stp	x29, x30, [sp, #-16-pbytes]!
	mov	x29, sp

	adr	x17, highctidh_1024_uintbig_p

	LD64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 0, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	subs
	ST64x8	sp, 16, x1, x2, x3, x4, x5, x6, x7, x8
	ST64x8	sp, 80, x9, x10, x11, x12, x13, x14, x15, x16

	LD64x8	x0, 64, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 64, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	sbcs

	cset	x17, cs
	neg	x17, x17

	CSWP2x8	x17
	ST64x8	x0, 64, x1, x2, x3, x4, x5, x6, x7, x8

	LD64x8	sp, 16, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	sp, 80, x9, x10, x11, x12, x13, x14, x15, x16
	CSWP2x8	x17
	ST64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8

	ldp	x29, x30, [sp], #16+pbytes
	ret


.global highctidh_1024_fp_add2
highctidh_1024_fp_add2:
	mov	x2, x0

.global highctidh_1024_fp_add3
highctidh_1024_fp_add3:
	stp	x29, x30, [sp, #-32]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_1024_uintbig_add3
	ldr	x0, [sp, #16]

	bl	reduce_once

	ldp	x29, x30, [sp], #32
	ret


.global highctidh_1024_fp_sub2
highctidh_1024_fp_sub2:
	mov	x2, x1
	mov	x1, x0

.global highctidh_1024_fp_sub3
highctidh_1024_fp_sub3:
	stp	x29, x30, [sp, #-16-pbytes]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_1024_uintbig_sub3
	neg	x1, x0
	ldr	x0, [sp, #16]

	adr	x2, highctidh_1024_uintbig_p
	.set k, 0
	.rept plimbs/2
		ldp	x3, x4, [x2, #8*k]
		and	x3, x3, x1
		and	x4, x4, x1
		stp	x3, x4, [sp, #16+8*k]
		.set k, k+2
	.endr

	ldp	x3, x4, [sp, #16]
	ldp	x5, x6, [x0]
	adds	x3, x3, x5
	adcs	x4, x4, x6
	stp	x3, x4, [x0]
	.set k, 2
	.rept (plimbs/2)-1
		ldp	x3, x4, [sp, #16+8*k]
		ldp	x5, x6, [x0, #8*k]
		adcs	x3, x3, x5
		adcs	x4, x4, x6
		stp	x3, x4, [x0, #8*k]
		.set k, k+2
	.endr

	ldp	x29, x30, [sp], #16+pbytes
	ret


.macro MULSTEP, k, I0, I1, I2, I3, I4, I5, I6, I7, I8, I9, I10, I11, I12, I13, I14, I15, I16
	ldr	xa_k, [x1, #8*\k]
	ldr	x5, [x2]

	mov	xcarry, xzr
	mov	xover, xzr

	LD64x8	xuintbig_p, 0, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I0, x10, \I1, x11, \I2, x12, \I3
	LDx4AT	sp, x13, \I4, x14, \I5, x15, \I6, x16, \I7

	mul	x3, xa_k, x5
	add	x3, x3, x9
	mul	x3, x3, xinv_min_p_mod_r

	MULX	x5, x4, x3, x21
	ADOX	x9, x4
	MAAx7	x3

	STx4AT	sp, x9, \I0, x10, \I1, x11, \I2, x12, \I3
	STx4AT	sp, x13, \I4, x14, \I5, x15, \I6, x16, \I7


	LD64x8	xuintbig_p, 64, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I8, x10, \I9, x11, \I10, x12, \I11
	LDx4AT	sp, x13, \I12, x14, \I13, x15, \I14, x16, \I15
	MAAx8	x3
	STx4AT	sp, x9, \I8, x10, \I9, x11, \I10, x12, \I11
	STx4AT	sp, x13, \I12, x14, \I13, x15, \I14, x16, \I15


	ldr	x9, [sp, #\I16]
	ADCX	x9, x5
	add	x9, x9, xover
	str	x9, [sp, #\I16]


	mov	xcarry, xzr
	mov	xover, xzr

	LD64x8	x2, 0, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I0, x10, \I1, x11, \I2, x12, \I3
	LDx4AT	sp, x13, \I4, x14, \I5, x15, \I6, x16, \I7

	MULX	x5, x4, xa_k, x21
	ADOX	x9, x4
	MAAx7	xa_k

	STx4AT	sp, x9, \I0, x10, \I1, x11, \I2, x12, \I3
	STx4AT	sp, x13, \I4, x14, \I5, x15, \I6, x16, \I7


	LD64x8	x2, 64, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I8, x10, \I9, x11, \I10, x12, \I11
	LDx4AT	sp, x13, \I12, x14, \I13, x15, \I14, x16, \I15
	MAAx8	xa_k
	STx4AT	sp, x9, \I8, x10, \I9, x11, \I10, x12, \I11
	STx4AT	sp, x13, \I12, x14, \I13, x15, \I14, x16, \I15


	ldr	x9, [sp, #\I16]
	ADCX	x9, x5
	add	x9, x9, xover
	str	x9, [sp, #\I16]
.endm

.global highctidh_1024_fp_mul2
highctidh_1024_fp_mul2:
	mov	x2, x0

.global highctidh_1024_fp_mul3
highctidh_1024_fp_mul3:

xinv_min_p_mod_r	.req	x7
xuintbig_p		.req	x8
xa_k			.req	x17
xcarry			.req	x19
xover			.req	x20

	stp	x29, x30, [sp, #-96-pbytes-16]!
	mov	x29, sp
	ST64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	stp	x27, x28, [sp, #80]

	adr	xinv_min_p_mod_r, .highctidh_1024_inv_min_p_mod_r
	ldr	xinv_min_p_mod_r, [xinv_min_p_mod_r]

	adr	xuintbig_p, highctidh_1024_uintbig_p

	add	sp, sp, #96

	mov	w5, #0x0001
	str	w5, [sp]
	ldr	p0, [sp]

	mov	w5, #0x0100
	str	w5, [sp]
	ldr	p1, [sp]

	ST64x8	sp, 0, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	ST64x8	sp, 64, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	str	xzr, [sp, #128]

	.set k, 0
	.rept plimbs
		MULSTEP	k, 8*(k+1), (8*(k+2))%136, (8*(k+3))%136, (8*(k+4))%136, (8*(k+5))%136, (8*(k+6))%136, (8*(k+7))%136, (8*(k+8))%136, (8*(k+9))%136, (8*(k+10))%136, (8*(k+11))%136, (8*(k+12))%136, (8*(k+13))%136, (8*(k+14))%136, (8*(k+15))%136, (8*(k+16))%136, 8*k
		.set k, k+1
	.endr

	mov	x1, sp
	sub	sp, sp, #96
	bl	highctidh_1024_fp_copy
	bl	reduce_once

	LD64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	ldp	x27, x28, [sp, #80]
	ldp	x29, x30, [sp], #96+pbytes+16
	ret


.global highctidh_1024_fp_sq1
highctidh_1024_fp_sq1:
	mov	x1, x0

.global highctidh_1024_fp_sq2
highctidh_1024_fp_sq2:
	mov	x2, x1
	b	highctidh_1024_fp_mul3
