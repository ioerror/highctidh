#include "aarch64_macros.s"

.set plimbs, 8
.set pbytes, plimbs*8

.section .rodata

#include "fp_const_le_512.s"

.section .text
.p2align 4,,15

.extern highctidh_512_uintbig_p
.hidden highctidh_512_uintbig_p


.global highctidh_512_fp_copy
highctidh_512_fp_copy:
	LD64x8	x1, 0, x2, x3, x4, x5, x6, x7, x8, x9
	ST64x8	x0, 0, x2, x3, x4, x5, x6, x7, x8, x9
	ret


.global highctidh_512_fp_cmov
highctidh_512_fp_cmov:
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


.global highctidh_512_fp_cswap
highctidh_512_fp_cswap:
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
	LD64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8

	adr	x17, highctidh_512_uintbig_p
	LD64x8	x17, 0, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	subs

	cset	x17, cs
	neg	x17, x17

	CSWP2x8	x17
	ST64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8
	ret


.global highctidh_512_fp_add2
highctidh_512_fp_add2:
	mov	x2, x0

.global highctidh_512_fp_add3
highctidh_512_fp_add3:
	stp	x29, x30, [sp, #-32]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_512_uintbig_add3
	ldr	x0, [sp, #16]

	bl	reduce_once
	ldp	x29, x30, [sp], #32
	ret


.global highctidh_512_fp_sub2
highctidh_512_fp_sub2:
	mov	x2, x1
	mov	x1, x0

.global highctidh_512_fp_sub3
highctidh_512_fp_sub3:
	stp	x29, x30, [sp, #-32]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_512_uintbig_sub3
	neg	x1, x0
	ldr	x0, [sp, #16]

	sub	sp, sp, pbytes

	adr	x2, highctidh_512_uintbig_p
	.set k, 0
	.rept plimbs/2
		ldp	x3, x4, [x2, #8*k]
		and	x3, x3, x1
		and	x4, x4, x1
		stp	x3, x4, [sp, #8*k]
		.set k, k+2
	.endr

	ldp	x3, x4, [sp]
	ldp	x5, x6, [x0]
	adds	x3, x3, x5
	adcs	x4, x4, x6
	stp	x3, x4, [x0]
	.set k, 2
	.rept (plimbs/2)-1
		ldp	x3, x4, [sp, #8*k]
		ldp	x5, x6, [x0, #8*k]
		adcs	x3, x3, x5
		adcs	x4, x4, x6
		stp	x3, x4, [x0, #8*k]
		.set k, k+2
	.endr

	add	sp, sp, pbytes
	ldp	x29, x30, [sp], #32
	ret


.macro MULSTEP, k, r0, r1, r2, r3, r4, r5, r6, r7, r8
	ldr	xa_k, [x1, #8*\k]
	ldr	x5, [x2]

	mul	x3, xa_k, x5
	add	x3, x3, \r0
	mul	x3, x3, xinv_min_p_mod_r

	mov	xcarry, xzr
	mov	xover, xzr

	adr	x5, highctidh_512_uintbig_p
	LD64x8	x5, 0, x21, x22, x23, x24, x25, x26, x27, x28

	MULX	x5, x4, x3, x21
	ADOX	\r0, x4

	ADCX	\r1, x5
	MULX	x5, x4, x3, x22
	ADOX	\r1, x4

	ADCX	\r2, x5
	MULX	x5, x4, x3, x23
	ADOX	\r2, x4

	ADCX	\r3, x5
	MULX	x5, x4, x3, x24
	ADOX	\r3, x4

	ADCX	\r4, x5
	MULX	x5, x4, x3, x25
	ADOX	\r4, x4

	ADCX	\r5, x5
	MULX	x5, x4, x3, x26
	ADOX	\r5, x4

	ADCX	\r6, x5
	MULX	x5, x4, x3, x27
	ADOX	\r6, x4

	ADCX	\r7, x5
	MULX	x5, x4, x3, x28
	ADOX	\r7, x4

	ADCX	\r8, x5
	add	\r8, \r8, xover


	mov	xcarry, xzr
	mov	xover, xzr

	LD64x8	x2, 0, x21, x22, x23, x24, x25, x26, x27, x28

	MULX	x5, x4, xa_k, x21
	ADOX	\r0, x4

	ADCX	\r1, x5
	MULX	x5, x4, xa_k, x22
	ADOX	\r1, x4

	ADCX	\r2, x5
	MULX	x5, x4, xa_k, x23
	ADOX	\r2, x4

	ADCX	\r3, x5
	MULX	x5, x4, xa_k, x24
	ADOX	\r3, x4

	ADCX	\r4, x5
	MULX	x5, x4, xa_k, x25
	ADOX	\r4, x4

	ADCX	\r5, x5
	MULX	x5, x4, xa_k, x26
	ADOX	\r5, x4

	ADCX	\r6, x5
	MULX	x5, x4, xa_k, x27
	ADOX	\r6, x4

	ADCX	\r7, x5
	MULX	x5, x4, xa_k, x28
	ADOX	\r7, x4

	ADCX	\r8, x5
	add	\r8, \r8, xover
.endm

.global highctidh_512_fp_mul2
highctidh_512_fp_mul2:
	mov	x2, x0

.global highctidh_512_fp_mul3
highctidh_512_fp_mul3:

xinv_min_p_mod_r	.req	x7
xa_k			.req	x17
xcarry			.req	x19
xover			.req	x20

	stp	x29, x30, [sp, #-96]!
	mov	x29, sp
	ST64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	stp	x27, x28, [sp, #80]

	adr	xinv_min_p_mod_r, .highctidh_512_inv_min_p_mod_r
	ldr	xinv_min_p_mod_r, [xinv_min_p_mod_r]

	mov	x8, xzr
	mov	x9, xzr
	mov	x10, xzr
	mov	x11, xzr
	mov	x12, xzr
	mov	x13, xzr
	mov	x14, xzr
	mov	x15, xzr
	mov	x16, xzr

	MULSTEP	0, x8,  x9,  x10, x11, x12, x13, x14, x15, x16
	MULSTEP	1, x9,  x10, x11, x12, x13, x14, x15, x16, x8
	MULSTEP	2, x10, x11, x12, x13, x14, x15, x16, x8,  x9
	MULSTEP	3, x11, x12, x13, x14, x15, x16, x8,  x9,  x10
	MULSTEP	4, x12, x13, x14, x15, x16, x8,  x9,  x10, x11
	MULSTEP	5, x13, x14, x15, x16, x8,  x9,  x10, x11, x12
	MULSTEP	6, x14, x15, x16, x8,  x9,  x10, x11, x12, x13
	MULSTEP	7, x15, x16, x8,  x9,  x10, x11, x12, x13, x14

	ST64x8	x0, 0, x16, x8, x9, x10, x11, x12, x13, x14

	bl	reduce_once

	LD64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	ldp	x27, x28, [sp, #80]
	ldp	x29, x30, [sp], #96
	ret


.global highctidh_512_fp_sq1
highctidh_512_fp_sq1:
	mov	x1, x0

.global highctidh_512_fp_sq2
highctidh_512_fp_sq2:
	mov	x2, x1
	b	highctidh_512_fp_mul3
