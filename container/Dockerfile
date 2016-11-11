FROM ubuntu:xenial
MAINTAINER Gabe Fierro <gtfierro@eecs.berkeley.edu>

RUN apt-get update && apt-get install -y libssl-dev graphviz libcurl3 libcurl3-gnutls libyajl-dev libxslt1-dev && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY server/ /server
ADD libraptor2.so.0 /usr/lib/.
ADD entrypoint.sh /bin/.

ENTRYPOINT [ "bin/entrypoint.sh" ]
