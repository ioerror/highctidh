#include "aarch64_macros.s"

.set plimbs, 32
.set pbytes, plimbs*8

.section .rodata

#include "fp_const_le_2048.s"

.section .text
.p2align 4,,15

.extern highctidh_2048_uintbig_p
.hidden highctidh_2048_uintbig_p


.global highctidh_2048_fp_copy
highctidh_2048_fp_copy:
	LD64x8	x1, 0, x2, x3, x4, x5, x6, x7, x8, x9
	LD64x8	x1, 64, x10, x11, x12, x13, x14, x15, x16, x17
	ST64x8	x0, 0, x2, x3, x4, x5, x6, x7, x8, x9
	ST64x8	x0, 64, x10, x11, x12, x13, x14, x15, x16, x17
	LD64x8	x1, 128, x2, x3, x4, x5, x6, x7, x8, x9
	LD64x8	x1, 192, x10, x11, x12, x13, x14, x15, x16, x17
	ST64x8	x0, 128, x2, x3, x4, x5, x6, x7, x8, x9
	ST64x8	x0, 192, x10, x11, x12, x13, x14, x15, x16, x17
	ret


.global highctidh_2048_fp_cmov
highctidh_2048_fp_cmov:
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


.global highctidh_2048_fp_cswap
highctidh_2048_fp_cswap:
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
	stp	x29, x30, [sp, #-16-384]!
	mov	x29, sp

	adr	x17, highctidh_2048_uintbig_p

	LD64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 0, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	subs
	ST64x8	sp, 16, x1, x2, x3, x4, x5, x6, x7, x8
	ST64x8	sp, 80, x9, x10, x11, x12, x13, x14, x15, x16

	LD64x8	x0, 64, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 64, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	sbcs
	ST64x8	sp, 144, x1, x2, x3, x4, x5, x6, x7, x8
	ST64x8	sp, 208, x9, x10, x11, x12, x13, x14, x15, x16

	LD64x8	x0, 128, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 128, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	sbcs
	ST64x8	sp, 272, x1, x2, x3, x4, x5, x6, x7, x8
	ST64x8	sp, 336, x9, x10, x11, x12, x13, x14, x15, x16

	LD64x8	x0, 192, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	x17, 192, x9, x10, x11, x12, x13, x14, x15, x16
	SBCSx8	sbcs

	cset	x17, cs
	neg	x17, x17

	CSWP2x8	x17
	ST64x8	x0, 192, x1, x2, x3, x4, x5, x6, x7, x8

	LD64x8	sp, 272, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	sp, 336, x9, x10, x11, x12, x13, x14, x15, x16
	CSWP2x8	x17
	ST64x8	x0, 128, x1, x2, x3, x4, x5, x6, x7, x8

	LD64x8	sp, 144, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	sp, 208, x9, x10, x11, x12, x13, x14, x15, x16
	CSWP2x8	x17
	ST64x8	x0, 64, x1, x2, x3, x4, x5, x6, x7, x8

	LD64x8	sp, 16, x1, x2, x3, x4, x5, x6, x7, x8
	LD64x8	sp, 80, x9, x10, x11, x12, x13, x14, x15, x16
	CSWP2x8	x17
	ST64x8	x0, 0, x1, x2, x3, x4, x5, x6, x7, x8

	ldp	x29, x30, [sp], #16+384
	ret


.global highctidh_2048_fp_add2
highctidh_2048_fp_add2:
	mov	x2, x0

.global highctidh_2048_fp_add3
highctidh_2048_fp_add3:
	stp	x29, x30, [sp, #-32]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_2048_uintbig_add3
	ldr	x0, [sp, #16]

	bl	reduce_once

	ldp	x29, x30, [sp], #32
	ret


.global highctidh_2048_fp_sub2
highctidh_2048_fp_sub2:
	mov	x2, x1
	mov	x1, x0

.global highctidh_2048_fp_sub3
highctidh_2048_fp_sub3:
	stp	x29, x30, [sp, #-16-pbytes]!
	mov	x29, sp

	str	x0, [sp, #16]
	bl	highctidh_2048_uintbig_sub3
	neg	x1, x0
	ldr	x0, [sp, #16]

	adr	x2, highctidh_2048_uintbig_p
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


.macro MULSTEP, k, I0, I1, I2, I3, I4, I5, I6, I7, I8, I9, I10, I11, I12, I13, I14, I15, I16, I17, I18, I19, I20, I21, I22, I23, I24, I25, I26, I27, I28, I29, I30, I31, I32
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


	LD64x8	xuintbig_p, 128, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I16, x10, \I17, x11, \I18, x12, \I19
	LDx4AT	sp, x13, \I20, x14, \I21, x15, \I22, x16, \I23
	MAAx8	x3
	STx4AT	sp, x9, \I16, x10, \I17, x11, \I18, x12, \I19
	STx4AT	sp, x13, \I20, x14, \I21, x15, \I22, x16, \I23


	LD64x8	xuintbig_p, 192, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I24, x10, \I25, x11, \I26, x12, \I27
	LDx4AT	sp, x13, \I28, x14, \I29, x15, \I30, x16, \I31
	MAAx8	x3
	STx4AT	sp, x9, \I24, x10, \I25, x11, \I26, x12, \I27
	STx4AT	sp, x13, \I28, x14, \I29, x15, \I30, x16, \I31


	ldr	x9, [sp, #\I32]
	ADCX	x9, x5
	add	x9, x9, xover
	str	x9, [sp, #\I32]


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


	LD64x8	x2, 128, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I16, x10, \I17, x11, \I18, x12, \I19
	LDx4AT	sp, x13, \I20, x14, \I21, x15, \I22, x16, \I23
	MAAx8	xa_k
	STx4AT	sp, x9, \I16, x10, \I17, x11, \I18, x12, \I19
	STx4AT	sp, x13, \I20, x14, \I21, x15, \I22, x16, \I23


	LD64x8	x2, 192, x21, x22, x23, x24, x25, x26, x27, x28
	LDx4AT	sp, x9, \I24, x10, \I25, x11, \I26, x12, \I27
	LDx4AT	sp, x13, \I28, x14, \I29, x15, \I30, x16, \I31
	MAAx8	xa_k
	STx4AT	sp, x9, \I24, x10, \I25, x11, \I26, x12, \I27
	STx4AT	sp, x13, \I28, x14, \I29, x15, \I30, x16, \I31


	ldr	x9, [sp, #\I32]
	ADCX	x9, x5
	add	x9, x9, xover
	str	x9, [sp, #\I32]
.endm

.global highctidh_2048_fp_mul2
highctidh_2048_fp_mul2:
	mov	x2, x0

.global highctidh_2048_fp_mul3
highctidh_2048_fp_mul3:

xinv_min_p_mod_r	.req	x7
xuintbig_p		.req	x8
xa_k			.req	x17
xcarry			.req	x19
xover			.req	x20

	stp	x29, x30, [sp, #-96-pbytes-16]!
	mov	x29, sp
	ST64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	stp	x27, x28, [sp, #80]

	adr	xinv_min_p_mod_r, .highctidh_2048_inv_min_p_mod_r
	ldr	xinv_min_p_mod_r, [xinv_min_p_mod_r]

	adr	xuintbig_p, highctidh_2048_uintbig_p

	add	sp, sp, #96

	ST64x8	sp, 0, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	ST64x8	sp, 64, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	ST64x8	sp, 128, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	ST64x8	sp, 192, xzr, xzr, xzr, xzr, xzr, xzr, xzr, xzr
	str	xzr, [sp, #256]

	MULSTEP	0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0
	MULSTEP	1, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8
	MULSTEP	2, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16
	MULSTEP	3, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24
	MULSTEP	4, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32
	MULSTEP	5, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40
	MULSTEP	6, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48
	MULSTEP	7, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56
	MULSTEP	8, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64
	MULSTEP	9, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72
	MULSTEP	10, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80
	MULSTEP	11, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88
	MULSTEP	12, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96
	MULSTEP	13, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104
	MULSTEP	14, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112
	MULSTEP	15, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120
	MULSTEP	16, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128
	MULSTEP	17, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136
	MULSTEP	18, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144
	MULSTEP	19, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152
	MULSTEP	20, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160
	MULSTEP	21, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168
	MULSTEP	22, 184, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176
	MULSTEP	23, 192, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184
	MULSTEP	24, 200, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192
	MULSTEP	25, 208, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200
	MULSTEP	26, 216, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208
	MULSTEP	27, 224, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216
	MULSTEP	28, 232, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224
	MULSTEP	29, 240, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232
	MULSTEP	30, 248, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240
	MULSTEP	31, 256, 0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192, 200, 208, 216, 224, 232, 240, 248

	mov	x1, sp
	sub	sp, sp, #96
	bl	highctidh_2048_fp_copy
	bl	reduce_once

	LD64x8	sp, 16, x19, x20, x21, x22, x23, x24, x25, x26
	ldp	x27, x28, [sp, #80]
	ldp	x29, x30, [sp], #96+pbytes+16
	ret


.global highctidh_2048_fp_sq1
highctidh_2048_fp_sq1:
	mov	x1, x0

.global highctidh_2048_fp_sq2
highctidh_2048_fp_sq2:
	mov	x2, x1
	b	highctidh_2048_fp_mul3
