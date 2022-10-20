# syntax=docker/dockerfile:1
ARG ARCH=
FROM ${ARCH}/debian:sid-slim
RUN set -eux; \
	apt update; \
	apt install -y --no-install-recommends make gcc clang 2>&1 > /dev/null; \
    apt clean; rm -rf /var/lib/apt/lists/*;
WORKDIR /highctidh/
COPY . /highctidh/
