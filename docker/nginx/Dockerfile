FROM base-rhel7

LABEL maintainer="Justin Cook <jhcook@secnix.com>"

RUN yum install -y nginx && \
    yum clean all && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log

VOLUME ["/var/cache/nginx"]

COPY ./hangman.conf /etc/nginx/conf.d/hangman.conf
COPY ./site /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
