#!/usr/bin/env bash
#
# This sets the appropriate environment variables, creates a CDK with access
# to appropriate insecure registries and allows transparent access.
#
# Author: Justin Cook <jhcook@secnix.com>

# Check if we have an environment that will work

export MINISHIFT_USERNAME="jhcook@secnix.com"
echo "Please enter your RHDS Password: "
read -sr MINISHIFT_PASSWORD_INPUT
export MINISHIFT_PASSWORD=$MINISHIFT_PASSWORD_INPUT

# minishift config set memory 8192

#minishift start --vm-driver xhyve
minishift start --vm-driver virtualbox

minishift ssh << __EOF__
sudo yum install -y git wget 
sudo yum clean all

wget https://github.com/openshift/source-to-image/releases/download/v1.1.6/source-to-image-v1.1.6-f519129-linux-amd64.tar.gz 
tar zxvf source-to-image-v1.1.6-f519129-linux-amd64.tar.gz ./s2i
sudo mv s2i /usr/bin/ 

#wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
#tar zxvf go1.8.3.linux-amd64.tar.gz 
rm -fr *
__EOF__
