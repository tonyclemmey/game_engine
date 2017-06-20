To build this container:
```
$ docker build --rm -t nginx-hangman .
```


To run this container:
```
$ docker run -d --rm -p 8080:8080 --name nginx-hangman nginx-hangman
```
