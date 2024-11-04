.macro SET, size
	stp	x1, xzr, [x0]
	.set k, 2
	.rept \size/64/2-1
		stp	xzr, xzr, [x0, #8*k]
		.set k, k+2
	.endr
	ret
.endm

.macro BIT, size
	and	x2, x1, #63
	and	x1, x1, #\size-64
	lsr	x1, x1, #3
	ldr	x0, [x0, x1]
	lsr	x0, x0, x2
	and	x0, x0, #1
	ret
.endm

.macro ADD3, size
	ldp	x3, x6, [x1]
	ldp	x4, x7, [x2]
	adds	x5, x3, x4
	adcs	x8, x6, x7
	stp	x5, x8, [x0]

	.set k, 2
	.rept \size/64/2-1
		ldp	x3, x6, [x1, #8*k]
		ldp	x4, x7, [x2, #8*k]
		adcs	x5, x3, x4
		adcs	x8, x6, x7
		stp	x5, x8, [x0, #8*k]
		.set k, k+2
	.endr

	cset	x0, cs
	ret
.endm

.macro SUB3, size
	ldp	x3, x6, [x1]
	ldp	x4, x7, [x2]
	subs	x5, x3, x4
	sbcs	x8, x6, x7
	stp	x5, x8, [x0]

	.set k, 2
	.rept \size/64/2-1
		ldp	x3, x6, [x1, #8*k]
		ldp	x4, x7, [x2, #8*k]
		sbcs	x5, x3, x4
		sbcs	x8, x6, x7
		stp	x5, x8, [x0, #8*k]
		.set k, k+2
	.endr

	cset	x0, cc
	ret
.endm

.macro MUL3_64, size
	ldp	x3, x6, [x1]
	mul	x4, x3, x2
	umulh	x5, x3, x2

	mul	x7, x6, x2
	adds	x7, x7, x5
	stp	x4, x7, [x0]

	.set k, 2
	.rept \size/64/2-1
		umulh	x5, x6, x2

		ldp	x3, x6, [x1, #8*k]

		mul	x4, x3, x2
		adcs	x4, x4, x5

		umulh	x5, x3, x2
		mul	x7, x6, x2
		adcs	x7, x7, x5

		stp	x4, x7, [x0, #8*k]
		.set k, k+2
	.endr
	ret
.endm

.macro LD64x8 rsrc, ofs, r1, r2, r3, r4, r5, r6, r7, r8
	ldp	\r1, \r2, [\rsrc, #\ofs]
	ldp	\r3, \r4, [\rsrc, #\ofs + 16]
	ldp	\r5, \r6, [\rsrc, #\ofs + 32]
	ldp	\r7, \r8, [\rsrc, #\ofs + 48]
.endm

.macro ST64x8 rdst, ofs, r1, r2, r3, r4, r5, r6, r7, r8
	stp	\r1, \r2, [\rdst, #\ofs]
	stp	\r3, \r4, [\rdst, #\ofs + 16]
	stp	\r5, \r6, [\rdst, #\ofs + 32]
	stp	\r7, \r8, [\rdst, #\ofs + 48]
.endm

.macro LDx4AT rsrc, r0, ofs0, r1, ofs1, r2, ofs2, r3, ofs3
	ldr	\r0, [\rsrc, #\ofs0]
	ldr	\r1, [\rsrc, #\ofs1]
	ldr	\r2, [\rsrc, #\ofs2]
	ldr	\r3, [\rsrc, #\ofs3]
.endm

.macro STx4AT rdst, r0, ofs0, r1, ofs1, r2, ofs2, r3, ofs3
	str	\r0, [\rdst, #\ofs0]
	str	\r1, [\rdst, #\ofs1]
	str	\r2, [\rdst, #\ofs2]
	str	\r3, [\rdst, #\ofs3]
.endm

.macro SBCSx8 insn1
	\insn1	x9, x1, x9
	sbcs	x10, x2, x10
	sbcs	x11, x3, x11
	sbcs	x12, x4, x12
	sbcs	x13, x5, x13
	sbcs	x14, x6, x14
	sbcs	x15, x7, x15
	sbcs	x16, x8, x16
.endm

.macro CSWP2, r0, r1, rmask
	eor	\r1, \r1, \r0
	and	\r1, \r1, \rmask
	eor	\r0, \r0, \r1
.endm

.macro CSWP2x8, rmask
	CSWP2	x1, x9, \rmask
	CSWP2	x2, x10, \rmask
	CSWP2	x3, x11, \rmask
	CSWP2	x4, x12, \rmask
	CSWP2	x5, x13, \rmask
	CSWP2	x6, x14, \rmask
	CSWP2	x7, x15, \rmask
	CSWP2	x8, x16, \rmask
.endm

.macro MULX roh, rol, ri1, ri2
	umulh	\roh, \ri1, \ri2
	mul	\rol, \ri1, \ri2
.endm

.macro ADCX r0, r1
	adds	\r0, \r0, xcarry
	cset	xcarry, cs
	adds	\r0, \r0, \r1
	adc	xcarry, xcarry, xzr
.endm

.macro ADOX r0, r1
	adds	\r0, \r0, xover
	cset	xover, cs
	adds	\r0, \r0, \r1
	adc	xover, xover, xzr
.endm
