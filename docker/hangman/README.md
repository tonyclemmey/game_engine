In this directory you create the base RHEL7 build as follows:
```
$ docker build --rm -t base-rhel7 .
...
Step 4 : COPY passwd.template /home/hangman/passwd.template
 ---> 4ad4614b84f6
Removing intermediate container ae50ff6096f9
Step 5 : WORKDIR /opt/app-root
 ---> Running in c97f3a9197a4
 ---> c8bef6e33e9f
Removing intermediate container c97f3a9197a4
Step 6 : USER 1001
 ---> Running in 98ef09f618cb
 ---> 4c9dff4665cd
Removing intermediate container 98ef09f618cb
Step 7 : CMD usage
 ---> Running in 39837d6bfed7
 ---> e14c1af9f903
Removing intermediate container 39837d6bfed7
Successfully built e14c1af9f903
```

You must then change to the root of the repo and build the image using s2i
builder and runtime images:

Create environment file `.s2i/environment`
```
__AK__=s0m3l0ngr4nd0m4pp11c4710nk3y
__AID__=appid
```

```
$ s2i build . golang:1.8 hangman --runtime-image base-rhel7 --runtime-artifact /tmp/src/hangman_engine -s file://../../.s2i/bin --environment-file ../../.s2i/environment
...
Installed:
  words.noarch 0:3.0-22.el7

Complete!
Loaded plugins: ovl, product-id, search-disabled-repos, subscription-manager
Cleaning repos: rhel-7-server-rpms rhel-7-server-rt-beta-rpms
              : rhel-7-server-rt-rpms
Cleaning up everything
Build completed successfully
```

The image can now be ran as follows:
```
$ docker run -p 3000:3000 --name hangman hangman
```
