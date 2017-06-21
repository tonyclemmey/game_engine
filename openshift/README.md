Building
==============
**Before you begin you need an app id and app key from Oxford Dictionary API**
https://developer.oxforddictionaries.com

Install the Red Hat CDK, and after you have an instance running you can access
the box via `minishift ssh` or use the `oc` cli that is dynamically linked in
your environment. There are more than one way to do this, but we use both
interfaces in this tutorial.

1. `minishift ssh` and clone this git repository.
2. Create the environment file.
3. Build the base Docker image.
4. Build the hangman image using a builder image and the base image you just
created.
5. Adjust the hangman.conf, hangman.js and hangman.html accordingly to match
URL hostnames.
6. Build the nginx-hangman image.
```
$ git clone https://github.com/jhcook/game_engine.git
$ cd game_engine
$ cat - > .s2i/environment << __EOF__
> __AK__=<replace with your app key>
> __AID__=<replace with your app id>
> __EOF__
$ cd docker/hangman
$ sudo docker build --rm -t base-rhel7 .
...
$ cd ../..
$ s2i build . golang:1.8 hangman --runtime-image base-rhel7 --runtime-artifact /tmp/src/hangman_engine -s file://.s2i/bin --environment-file .s2i/environment
...
$ cd docker/nginx
$ sudo docker build --rm -t nginx-hangman .
...
```

7. Using the `oc` client -- perhaps in another terminal, create a service 
account.
8. Grant the service account privileges 
9. Export a token for use
```
$ oc login -u system:admin
...
$ oc create serviceaccount pusher -n myproject
...
$ oc policy add-role-to-user system:image-builder system:serviceaccount:myproject:pusher
role "system:image-builder" added: "system:serviceaccount:myproject:pusher"
$ oc describe sa pusher -n myproject
Name:        pusher
Namespace:    myproject
Labels:        <none>

Image pull secrets:    pusher-dockercfg-69v6q

Mountable secrets:     pusher-dockercfg-69v6q
                       pusher-token-kd2q3

Tokens:                pusher-token-g7ssq
                       pusher-token-kd2q3
$ oc describe secret pusher-token-g7ssq -n myproject
Name:        pusher-token-g7ssq
Namespace:    myproject
Labels:        <none>
Annotations:    kubernetes.io/created-by=openshift.io/create-dockercfg-secrets
        kubernetes.io/service-account.name=pusher
        kubernetes.io/service-account.uid=e0bca293-55da-11e7-8a93-080027026ffd

Type:    kubernetes.io/service-account-token

Data
====
service-ca.crt:    2186 bytes
token:        eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoicHVzaGVyLXRva2VuLWc3c3NxIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InB1c2hlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUwYmNhMjkzLTU1ZGEtMTFlNy04YTkzLTA4MDAyNzAyNmZmZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpteXByb2plY3Q6cHVzaGVyIn0.QgqN3QXh484ykCwpSAEOOaOk5aZToHf2IOhE0TdWQRhunk5fjlOXhVAEpvhlhwZxMXn2QWsPxO2HISd96-Jy7MeWnlce7gMgzbw_8yvu52okfhlLABs8lOBLJCkJuS0dha51RtDb-COCDkog2nNISB71ucX0ZNU_8WLJsxQ_tHxyBER0nkcOz2A0tnB-TDUYMy4hhfKEyPGC00RPiLhm5DWikwvCzhmQ4GTEnHMp2qxnKkqGPIpaktYl1IjA9D7lOgNPV0YP97X7t6BWnsUXRruW7QR9Qf3H6aqLPNDNFuZvUPDONZOqwE1ngaqxU2rtCD9bG-W4bU2WyjNBLP-v4w
ca.crt:        1070 bytes
namespace:    9 bytes
$ oc get svc/docker-registry -n default
NAME              CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
docker-registry   172.30.1.1   <none>        5000/TCP   1h
```

Using the token listed in the secret, return to `minishift ssh` terminal, login
to Docker, tag the images, and push to the registry.
```
$ sudo docker login 172.30.1.1:5000 -u pusher -p eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoicHVzaGVyLXRva2VuLWc3c3NxIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InB1c2hlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUwYmNhMjkzLTU1ZGEtMTFlNy04YTkzLTA4MDAyNzAyNmZmZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpteXByb2plY3Q6cHVzaGVyIn0.QgqN3QXh484ykCwpSAEOOaOk5aZToHf2IOhE0TdWQRhunk5fjlOXhVAEpvhlhwZxMXn2QWsPxO2HISd96-Jy7MeWnlce7gMgzbw_8yvu52okfhlLABs8lOBLJCkJuS0dha51RtDb-COCDkog2nNISB71ucX0ZNU_8WLJsxQ_tHxyBER0nkcOz2A0tnB-TDUYMy4hhfKEyPGC00RPiLhm5DWikwvCzhmQ4GTEnHMp2qxnKkqGPIpaktYl1IjA9D7lOgNPV0YP97X7t6BWnsUXRruW7QR9Qf3H6aqLPNDNFuZvUPDONZOqwE1ngaqxU2rtCD9bG-W4bU2WyjNBLP-v4w
Login Succeeded
$ sudo docker tag nginx-hangman 172.30.1.1:5000/myproject/nginx-hangman
$ sudo docker tag hangman 172.30.1.1:5000/myproject/hangman
$ sudo docker push 172.30.1.1:5000/myproject/hangman
The push refers to a repository [172.30.1.1:5000/myproject/hangman]
bc5a59b6f00e: Pushed
299a04a9366f: Pushed
fb8582c7ffd0: Pushed
279bfd6c7049: Pushed
f5bd5357a1de: Pushed
latest: digest: sha256:e1e386aa3b6fa0776d0d2122357aef711cb437ff28e15af9eab8e334bc112d3d size: 10567
$ sudo docker tag nginx-hangman 172.30.1.1:5000/myproject/nginx-hangman
$ sudo docker tag hangman 172.30.1.1:5000/myproject/hangman
$ sudo docker push 172.30.1.1:5000/myproject/hangman
The push refers to a repository [172.30.1.1:5000/myproject/hangman]
3c7021900fa2: Pushed
7a000cc189a6: Pushed
9ff71f1fec1f: Pushed
279bfd6c7049: Pushed
f5bd5357a1de: Pushed
latest: digest: sha256:c2983327ee0422b8199fb773d9998e8ef4e81b1a03b1e1034e8756313435e302 size: 10373
$ sudo docker push 172.30.1.1:5000/myproject/nginx-hangman
The push refers to a repository [172.30.1.1:5000/myproject/nginx-hangman]
e00e055841e0: Pushed
f9c6cbec40e4: Pushed
a552ca691e49: Pushed
7487bf0353a7: Pushed
8781ec54ba04: Pushed
latest: digest: sha256:6bfbb2474196c5471fce67e779d7e97ffc4f866465387120cffde5d3e3c143ef size: 9885
```

You can now return to the other terminal and see the image streams:
```
$ oc get is
NAME            DOCKER REPO                               TAGS      UPDATED
hangman         172.30.1.1:5000/myproject/hangman         latest    About a minute ago
nginx-hangman   172.30.1.1:5000/myproject/nginx-hangman   latest    50 seconds ago
```

This is not a recommended setting, but for the purposes of this demonstration,
please allow the restricted scc to run as any user.
https://access.redhat.com/articles/1493353
```
$ oc edit scc restricted
Change the runAsUser.Type strategy to RunAsAny.
```

Create the deploymentconfig, the service and export with a route for use
externally.
```
$ oc create -f hangman-dc.yaml
$ oc create -f hangman-svc.yaml
$ oc create -f hangman-route.yaml
```
