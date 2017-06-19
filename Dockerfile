FROM rhel7

LABEL io.k8s.description="Platform for receiving artefacts from s2i" \
      io.openshift.s2i.scripts-url="image:///usr/local/s2i" \
      maintainer="Justin Cook jhcook@secnix.com"

COPY ./.s2i/bin/ /usr/local/s2i

RUN yum install -y --setopt=tsflags=nodocs sudo nss_wrapper gettext && \
     useradd -u 1001 hangman && echo "hangman ALL = (ALL) NOPASSWD: ALL" > \
    /etc/sudoers.d/hangman_conf && chmod 0440 /etc/sudoers.d/hangman_conf && \
    mkdir -p /opt/app-root && chgrp -R 0 /opt/app-root && chmod -R g+rwX \
    /opt/app-root && yum clean all

COPY passwd.template /opt/app-root/passwd.template

WORKDIR /opt/app-root

USER 1001

CMD ["usage"]
