FROM rhel7

LABEL io.k8s.description="Platform for receiving artefacts from s2i" \
      io.openshift.s2i.scripts-url="image:///usr/local/s2i" \
      maintainer="Justin Cook jhcook@secnix.com"

COPY ./.s2i/bin/ /usr/local/s2i

RUN yum-config-manager --disable rhel-7-server-rt-rpms && \
    yum-config-manager --disable rhel-7-server-rt-beta-rpms && \
    yum install -y --enablerepo rhel-server-rhscl-7-rpms \
    --setopt=tsflags=nodocs sudo nss_wrapper gettext && \
    useradd -u 1001 hangman && \
    echo "hangman ALL = (ALL) NOPASSWD: ALL" > /etc/sudoers.d/hangman_conf && \
    chmod 0440 /etc/sudoers.d/hangman_conf && \
    mkdir -p /opt/app-root && \
    chgrp -R 0 /opt/app-root && \
    chmod -R g+rwX /opt/app-root && \
    yum clean all

COPY passwd.template /home/hangman/passwd.template

WORKDIR /opt/app-root

USER 1001

CMD ["usage"]
