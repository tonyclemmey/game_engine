Building
==============
**Before you begin, you need an app id and app key from Oxford Dictionary API**
https://developer.oxforddictionaries.com

**Before you begin, you need a Red Hat Developer Subscription**
https://developers.redhat.com/articles/frequently-asked-questions-no-cost-red-hat-enterprise-linux-developer-suite/

Install the Red Hat CDK, and after you have an instance running you can access
the box via `minishift ssh` or use the `oc` cli that is dynamically linked in
your environment. There are more than one way to do this, but we use both
interfaces in this tutorial.

https://access.redhat.com/documentation/en-us/red_hat_container_development_kit/3.0/html/installation_guide/
https://www.youtube.com/watch?v=UxwBB0_-9VM

You will need Vagrant installed in order to take advantage of the Vagrantfile
in this tutorial. After you have installed Vagrant, you will need to install
the vagrant-registration plugin for use with your Red Hat Developer Subscription
credentials. `vagrant plugin install vagrant-registration`

You are now ready to bring up the development environment for this tutorial. 
Change to this directory and bring up the CDK. `./minishift_boot.sh` Please follow 
the prompts and you will be left with a usable CDK.

1.  Bring up the CDK from this directory in the repository.
2.  Export the admin credentials.
3.  Change to the root directory of this repository and issue `vagrant up`.
4.  `vagrant ssh`, import OpenShift config, `cd /vagrant`, and create the environment file.
5.  Build the base Docker image.
6.  Build the hangman image using a builder image and the base image you just
    created.
7.  Build the nginx-hangman image.

```
$ ./minishift_boot.sh
...
$ minishift ip
192.168.99.100
$ minishift ssh "sudo cat /mnt/sda1/var/lib/minishift/openshift.local.config/master/admin.kubeconfig" > \
  admin.kubeconfig
$ cd ..
$ vagrant up
...
$ vagrant ssh
[vagrant@rhel7-vbox ~]$ mkdir .kube
[vagrant@rhel7-vbox ~]$ cp /vagrant/openshift/admin.kubeconfig ~/.kube/config
[vagrant@rhel7-vbox ~]$ cd /vagrant
[vagrant@rhel7-vbox vagrant]$ cat - > .s2i/environment << __EOF__
> __AK__=<replace with your app key>
> __AID__=<replace with your app id>
> __EOF__
[vagrant@rhel7-vbox vagrant]$ cd docker/hangman
[vagrant@rhel7-vbox hangman]$ sudo docker build --rm -t base-rhel7 .
...
[vagrant@rhel7-vbox hangman]$ cd ../..
[vagrant@rhel7-vbox vagrant]$ sudo s2i build . golang:1.8 hangman --runtime-image base-rhel7 --runtime-artifact \
/tmp/src/hangman_engine -s file://.s2i/bin --environment-file .s2i/environment
...
[vagrant@rhel7-vbox vagrant]$ cd docker/nginx
[vagrant@rhel7-vbox nginx]$ sudo docker build --rm -t nginx-hangman .
...
[vagrant@rhel7-vbox nginx]$ cd
```

8.  Using the `oc` client, create a service account.
9.  Grant the service account privileges to push images
10. Export a token for use to login with `docker login`

```
[vagrant@rhel7-vbox ~]$ oc create serviceaccount pusher -n myproject
...
[vagrant@rhel7-vbox ~]$ oc policy add-role-to-user system:image-builder system:serviceaccount:myproject:pusher -n myproject
role "system:image-builder" added: "system:serviceaccount:myproject:pusher"
[vagrant@rhel7-vbox ~]$ oc describe sa pusher -n myproject
Name:        pusher
Namespace:    myproject
Labels:        <none>

Image pull secrets:    pusher-dockercfg-69v6q

Mountable secrets:     pusher-dockercfg-69v6q
                       pusher-token-kd2q3

Tokens:                pusher-token-g7ssq
                       pusher-token-kd2q3
[vagrant@rhel7-vbox ~]$ oc describe secret pusher-token-g7ssq -n myproject
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
token:        eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoicHVzaGVyLXRva2VuLTdtMTltIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InB1c2hlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI4NWJhMjYwLTVlNWMtMTFlNy1hYTVlLTA4MDAyNzhhMzI0OSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpteXByb2plY3Q6cHVzaGVyIn0.Es_PuSMoeteFD4oodN22UI7e9VBJbFLJsROye4aGxA2dn5_0glJkozYc92cIVwIZXdLFblSloB23rzllOf3_NCKC2wUcIvPkM1iKwuDRmz7n6wXUR3TzDJKMN8gu6lFwGY5WQlfAJDe9Vvfv8XnGnok_hPQ-os1pNQHPMKe5TCkb8RmwkTvAbeczK4QDlKhaYImdKFSMG8x3JbCpNRWY4N4I56tvgdPWzwCgFrfKH732tQbAW058e5vv0kzMHDw-so-2qE-nkZlqKu5zMGIX46jYUbffLumAqtu0rDoe9dDPmc_QkRpCLcMJJJ9HW_fp0kwDhIfQD0Fqrl60ONlGbQ
ca.crt:        1070 bytes
namespace:    9 bytes
[vagrant@rhel7-vbox ~]$ oc expose svc/docker-registry -n default
route "docker-registry" exposed
[vagrant@rhel7-vbox ~]$ oc get route/docker-registry -n default
NAME              HOST/PORT                                       PATH      SERVICES          PORT       TERMINATION   WILDCARD
docker-registry   docker-registry-default.192.168.99.100.nip.io             docker-registry   5000-tcp                 None
```

11. Create relevant /etc/hosts file entries for the registry and app.
12. Mark the registry as insecure so it will not be rejected. 
13. Login to Docker
14. Tag the images
15. Push to the OpenShift Docker registry

```
[vagrant@rhel7-vbox ~]$ sudo bash -c "echo 192.168.99.100 hangman.example.com >> /etc/hosts"
[vagrant@rhel7-vbox ~]$ sudo bash -c "echo 192.168.99.100 docker-registry-default.192.168.99.100.nip.io >> /etc/hosts"
[vagrant@rhel7-vbox ~]$ sudo sed -i "s|#\ INSECURE_REGISTRY=.*|INSECURE_REGISTRY='--insecure-registry docker-registry-default.192.168.99.100.nip.io:80'|" /etc/sysconfig/docker
[vagrant@rhel7-vbox ~]$ sudo systemctl restart docker
[vagrant@rhel7-vbox ~]$ sudo docker login docker-registry-default.192.168.99.100.nip.io:80 -u pusher -p eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoicHVzaGVyLXRva2VuLTdtMTltIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InB1c2hlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI4NWJhMjYwLTVlNWMtMTFlNy1hYTVlLTA4MDAyNzhhMzI0OSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpteXByb2plY3Q6cHVzaGVyIn0.Es_PuSMoeteFD4oodN22UI7e9VBJbFLJsROye4aGxA2dn5_0glJkozYc92cIVwIZXdLFblSloB23rzllOf3_NCKC2wUcIvPkM1iKwuDRmz7n6wXUR3TzDJKMN8gu6lFwGY5WQlfAJDe9Vvfv8XnGnok_hPQ-os1pNQHPMKe5TCkb8RmwkTvAbeczK4QDlKhaYImdKFSMG8x3JbCpNRWY4N4I56tvgdPWzwCgFrfKH732tQbAW058e5vv0kzMHDw-so-2qE-nkZlqKu5zMGIX46jYUbffLumAqtu0rDoe9dDPmc_QkRpCLcMJJJ9HW_fp0kwDhIfQD0Fqrl60ONlGbQ
Login Succeeded
[vagrant@rhel7-vbox ~]$ sudo docker tag hangman docker-registry-default.192.168.99.100.nip.io:80/myproject/hangman
[vagrant@rhel7-vbox ~]$ sudo docker tag nginx-hangman docker-registry-default.192.168.99.100.nip.io:80/myproject/nginx-hangman
[vagrant@rhel7-vbox ~]$ sudo docker push docker-registry-default.192.168.99.100.nip.io:80/myproject/hangman
The push refers to a repository [docker-registry-default.192.168.99.100.nip.io:80/myproject/hangman]
c3ae2c8497c4: Pushed
5a43b655a19d: Pushed
ce6d7c376e5b: Pushed
86888f0aea6d: Pushed
dda6e8dfdcf7: Pushed
latest: digest: sha256:e7f633d0b2175b5ea2b841a8cd1337df58d47413beedcbe39f1aa2128b4e67c5 size: 9646
[vagrant@rhel7-vbox ~]$ sudo docker push docker-registry-default.192.168.99.100.nip.io:80/myproject/nginx-hangman
The push refers to a repository [docker-registry-default.192.168.99.100.nip.io:80/myproject/nginx-hangman]
eccfc2e529b9: Pushed
fbe8eaa10c08: Pushed
1ddac1491315: Pushed
5a43b655a19d: Mounted from myproject/hangman
ce6d7c376e5b: Mounted from myproject/hangman
86888f0aea6d: Mounted from myproject/hangman
dda6e8dfdcf7: Mounted from myproject/hangman
latest: digest: sha256:a07f3807f9400b1fe7388035a3bf552c046726dd7ab1da2023495445e6f6f928 size: 12295
```

You can now return to the other terminal and see the image streams:

```
[vagrant@rhel7-vbox ~]$ oc get is
NAME            DOCKER REPO                               TAGS      UPDATED
hangman         172.30.1.1:5000/myproject/hangman         latest    About a minute ago
nginx-hangman   172.30.1.1:5000/myproject/nginx-hangman   latest    50 seconds ago
```

This is not a recommended setting, but for the purposes of this demonstration,
please allow the restricted scc to run as any user.
https://access.redhat.com/articles/1493353

```
[vagrant@rhel7-vbox ~]$ oc edit scc restricted
Change the runAsUser.Type strategy to RunAsAny.
```

Create the deploymentconfig, the service and export with a route for use
externally.

```
[vagrant@rhel7-vbox ~]$ cd /vagrant/openshift
[vagrant@rhel7-vbox openshift]$ oc create -f hangman-dc.yaml -n myproject
[vagrant@rhel7-vbox openshift]$ oc create -f hangman-svc.yaml -n myproject
[vagrant@rhel7-vbox openshift]$ oc create -f hangman-route.yaml -n myproject
```
