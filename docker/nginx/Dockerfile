FROM nginx

LABEL maintainer="Justin Cook <jhcook@secnix.com>"

COPY ./docker/nginx/hangman.conf /etc/nginx/conf.d/hangman.conf
COPY ./site /usr/share/nginx/html

EXPOSE 8080
