#undef NAMESPACEBITS
#undef FIAT_BITS
#undef primes_num
#undef primes_batches
#undef primes_maxbatchboundplussize
// git grep ifndef|grep -i _H|cut -d ' ' -f 2-|sort -u|sed 's,^,#undef ,g'
#undef ANNOTATIONS_H
#undef cpucycles_h
#undef crypto_classify_h
#undef crypto_declassify_h
#undef CSIDH_H
#undef ELLIGATOR_H
#undef fiat_p1024_H
#undef fiat_p2048_H
#undef fiat_p511_H
#undef fiat_p512_H
#undef FP_H
#undef int32mask_h
#undef int32_sort_H
#undef int64mask_h
#undef MONT_H
#undef POLY_H
#undef primes_h
#undef proj_h
#undef randombytes_h
#undef RANDOM_H
#undef STEPS_H
#undef uintbig_h

/* we keep BITS and define the new macros: */
#define NAMESPACEBITS3(a,x,y) a ## y ## _ ## x
#define NAMESPACEBITS2(a,x,y) NAMESPACEBITS3(a,x,y)
#define NAMESPACEBITS(x) NAMESPACEBITS2(highctidh_,x,BITS)
#define public_key NAMESPACEBITS(public_key)
#define private_key NAMESPACEBITS(private_key)
#define csidh NAMESPACEBITS(csidh)
#define csidh_private NAMESPACEBITS(csidh_private)
#define base NAMESPACEBITS(base)
#define public_key_from_bytes NAMESPACEBITS(public_key_from_bytes)
#define public_key_to_bytes NAMESPACEBITS(public_key_to_bytes)
#define csidh_private_withrng NAMESPACEBITS(csidh_private_withrng)
#define validate NAMESPACEBITS(validate)
#define validate_cutofforder_v2 NAMESPACEBITS(validate_cutofforder_v2)
#define action NAMESPACEBITS(action)
