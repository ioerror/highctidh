#!/usr/bin/env bash

set -e -u -o pipefail
set -x

# on i7-9750H this takes under a minute
ulimit -s unlimited && \
  time /fiat-crypto/src/ExtractionOCaml/word_by_word_montgomery \
  --widen-bytes --internal-static --lang C \
  -o /docker_build_output/fiat-crypto-generated-code/fiat_p512.c \
  p512 64 5326738796327623094747867617954605554069371494832722337612446642054009560026576537626892113026381253624626941643949444792662881241621373288942880288065659

# on i7-9750H this takes ~20 minutes
ulimit -s unlimited && \
  time /fiat-crypto/src/ExtractionOCaml/word_by_word_montgomery \
  --widen-bytes --internal-static --lang C \
  -o /docker_build_output/fiat-crypto-generated-code/fiat_p1024.c \
  p1024 64 10397125823368453045280646945602587680645373501407971524723936453540081402468763489237207685909470852314300607054838964876369377755966328162572638385477152807267330439309520268717606043270205686916407950751430775965737505772592082230763227336919124289210653129066316362633007647060631513541162047451995268179

# This run takes >128GB of RAM and may take over a solid day depending on the CPU
ulimit -s unlimited && \
  time /fiat-crypto/src/ExtractionOCaml/word_by_word_montgomery \
  --widen-bytes --internal-static --lang C \
  -o /docker_build_output/fiat-crypto-generated-code/fiat_p2048.c \
  p2048 64 8528475943061024837635249434131089507352961229039071418359007853431472824232731281573875097585132124495905689769457289097290485066318818670998665693553484375832840901492721916443754108187759409880104977491168497435797029889150028375522118007040227895709956565527693986603274933373650307807559167751481868507471338172008086897049878117533316735319086656307962660962929231692405440348962572770958110754184132520227272603171250943083809093992306118654477172531039336755846450097401488026516556186294026830749211238807256327438743073467304920067034434209788101583120190178099235423322734627165016835534905946794432350531
