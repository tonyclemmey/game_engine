FROM base-rhel7

LABEL maintainer="Justin Cook <jhcook@secnix.com>"

RUN sudo rpm -Uvh --nodeps http://nginx.org/packages/rhel/7/x86_64/RPMS/nginx-1.12.0-1.el7.ngx.x86_64.rpm && \
    sudo mkdir -p /var/cache/nginx/client_temp && \
    sudo chmod g+rwX /var/cache/nginx/client_temp &&\
    sudo rm -f /etc/nginx/conf.d/default.conf && \
    sudo ln -sf /dev/stdout /var/log/nginx/access.log && \
    sudo ln -sf /dev/stderr /var/log/nginx/error.log

COPY ./hangman.conf /etc/nginx/conf.d/hangman.conf
COPY ./site /usr/share/nginx/html

USER 0

CMD ["nginx", "-g", "daemon off;"]
