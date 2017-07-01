#!/usr/bin/env bash
#
# This sets the appropriate environment variables, creates a CDK with access
# to appropriate insecure registries, s2i, golang and allows transparent 
# access.
#
# Workarounds for exhausting live-rw filesystem
# https://github.com/minishift/minishift/issues/1076
#
# Author: Justin Cook <jhcook@secnix.com>

# Link to S2I
export S2ILINK="https://github.com/openshift/source-to-image/releases/download/v1.1.7/source-to-image-v1.1.7-226afa1-linux-amd64.tar.gz"

# Link to Go
export GOLINK="https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz"

# Change this to your Red Hat account login
export MINISHIFT_USERNAME="jhcook@secnix.com"
echo "Please enter your RHDS Password: "
read -sr MINISHIFT_PASSWORD_INPUT
export MINISHIFT_PASSWORD=$MINISHIFT_PASSWORD_INPUT

# If you want you can change the memory allocated
# minishift config set memory 8192

# You can use the xhyve driver and insecure registries as follows
#minishift start --vm-driver xhyve --insecure-registry 10.0.0.1
minishift start --vm-driver virtualbox

minishift ssh << __EOF__
# Backoff and try a couple times if yum install fails
counter=0
until [ "\$counter" -ge "3" ]
do
  sudo yum install --disablerepo rhel-7-server-rt-beta-rpms -y git wget && break
  echo "yum install failed. Waiting..."
  sleep 5
  let "counter++"
done

# Clean yum cache
sudo yum clean all

# Create /opt on the virtual drive as / is a dm with limited space
sudo mkdir -p /mnt/sda1/opt 
sudo rm -rf /opt 
sudo ln -s /mnt/sda1/opt /opt 
sudo chmod o+rwx /opt

# Change to /tmp, download packages, extract and clean up since we have limited
# dm space.
cd /tmp
wget $S2ILINK > /tmp/s2iinstall.out
sudo tar zxvf `basename $S2ILINK` -C /usr/bin/ ./s2i >> /tmp/s2iinstall.out
wget $GOLINK > /tmp/goinstall.out
sudo tar zxvf `basename $GOLINK` -C /opt >> /tmp/goinstall.out

# Clean up after ourselves
rm -fr * 2>/dev/null || /bin/true

# Update .bash_profile with relevant information
cd
# Nested here doc is the cool
cat - > ~/.bash_profile << __EOF2__
export GOROOT=/opt/go/
export PATH=\$PATH:/opt/go/bin 
__EOF2__
__EOF__
