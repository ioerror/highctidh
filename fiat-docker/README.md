This directory contains the pristine generated `fiat_*.c` files. Regeneration
is possible using the included `Dockerfile` and may be updated as the
`fiat-crypto` generated code improves.

The `fiat-crypto` synthesis is described in git commit
`0ab12b8ce4b5cfd8ed615679cd3a3bb768b481a7` and later commits. The main commit
describes the steps and software versions used to generate the fiat backend.
To re-generate the fiat-crypto C code there is basic automation to assist with
validation and verification:
```
  ./misc/docker-setup-fiat.sh
  ./misc/docker-generate-fiat.sh
```

The submodules for `fiat-crypto` were in the following state:
```
$ git submodule status
 762c7f5e8c5fa13c85020d5b24b7910843266f00 coqprime (v8.14.1-49-g762c7f5)
 1f0568fc9bead08425463fcbb9d086f6d0e7f21f etc/coq-scripts (heads/master)
 e62ff4f3a823e318f746b606009c963b05282569 rewriter (v0.0.7-2-ge62ff4f3a)
 ea154e322c698c9699241a185c10b026ea96b5ee rupicola (v0.0.6)
```
    
Furthermore, we used `opam switch create 4.13.1` for our build
environment. Coq was version `8.16.0`:

A patch is generated for each field size as follows:
```
diff fiat_p512.c ../fiat_p512.c > fiat_p512.patch
diff fiat_p1024.c ../fiat_p1024.c > fiat_p1024.patch
diff fiat_p2048.c ../fiat_p2048.c > fiat_p2048.patch
```
