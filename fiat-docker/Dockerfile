# syntax=docker/dockerfile:1
FROM debian:buster-slim
RUN set -eux; \
	apt update; \
    apt install -y --no-install-recommends ca-certificates git coq \
    ocaml-findlib libcoq-ocaml-dev jq opam m4 lsb-release bsdextrautils; \
    apt clean; rm -rf /var/lib/apt/lists/*;
RUN apt update; apt install -y libgmp-dev;
RUN opam init --disable-sandboxing -y;
RUN eval $(opam env) opam switch create 4.13.1;
RUN eval $(opam env) opam update;
RUN eval $(opam env) opam pin -y add coq 8.16.0;
RUN eval $(opam env) opam install -y coq;
RUN git clone https://github.com/mit-plv/fiat-crypto && \
    cd fiat-crypto/ && \
    git checkout fe63f5595406ac9b7edd12fe05c75ee276dd6f63 && \
    git submodule update --init --recursive
WORKDIR fiat-crypto/
RUN make clean
RUN eval $(opam env) make standalone-ocaml
RUN pwd
WORKDIR /highctidh/
COPY . /highctidh/
