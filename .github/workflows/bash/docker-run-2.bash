docker run \
  -v $(pwd):/${{ github.workspace }} \
  -w ${{ github.workspace }} \
  -e "CC=${{ matrix.CC }}" \
  -e "HIGHCTIDH_PORTABLE=${{ matrix.PORTABLE }}" \
  -e "CFLAGS=-fno-lto" \
  --rm \
  ${{ matrix.DOCKER_ARCH_ARGS}} \
  /bin/sh \
  -c "./misc/install-alpine-deps.sh && ./misc/docker-gha-fixup.sh && make alpine-multi-arch"
