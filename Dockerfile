# syntax=docker/dockerfile:1
ARG ARCH=
FROM ${ARCH}/debian:sid-slim
RUN set -eux; \
	apt update; \
	DEBIAN_FRONTEND=noninteractive apt install -y --no-install-recommends \
        make gcc clang git \
        python3 python3-build python3-setuptools python3-stdeb \
        build-essential python3-venv python3-wheel python3-pip \
        dh-python python3-all-dev flit fakeroot coreutils \
        python3-pytest python3-pytest-xdist \
        2>&1 > /dev/null; \
    apt clean; rm -rf /var/lib/apt/lists/*;
WORKDIR /highctidh/
