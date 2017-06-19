FROM rhel7

LABEL io.k8s.description="Platform for receiving artefacts from s2i" \
      io.openshift.s2i.scripts-url="image:///usr/local/s2i" \
      maintainer="Justin Cook jhcook@secnix.com"

COPY ./.s2i/bin/ /usr/local/s2i

RUN mkdir -p /opt/app-root ; chown -R 1001:1001 /opt/app-root
WORKDIR /opt/app-root

USER 1001

CMD ["usage"]
