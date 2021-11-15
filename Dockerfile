# *************************************
#
# OpenGM
#
# *************************************

FROM alpine:3.14

MAINTAINER XTech Cloud "xtech.cloud"

ENV container docker
ENV MSA_MODE release

EXPOSE 18804

ADD bin/ogm-license /usr/local/bin/
RUN chmod +x /usr/local/bin/ogm-license

CMD ["/usr/local/bin/ogm-license"]
