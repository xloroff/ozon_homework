FROM ubuntu/prometheus:latest AS development
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get -y install curl

FROM ubuntu/prometheus:latest AS production
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get -y install curl

FROM ubuntu/prometheus:latest AS test
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get -y install curl