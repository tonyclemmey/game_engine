# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "secnix/rhel7"
  config.vm.network "private_network", type: "dhcp",
                    virtualbox__intnet: "vboxnet1"
  config.vm.provision "shell", inline: <<-SHELL
    # Install and configure Docker
    subscription-manager repos --enable=rhel-7-server-extras-rpms
    subscription-manager repos --enable=rhel-7-server-optional-rpms
    yum install -y docker device-mapper-libs device-mapper-event-libs
    # Install Gnome and set as default target https://access.redhat.com/solutions/5238
    yum groupinstall -y gnome-desktop x11 fonts
    systemctl start docker.service
    systemctl enable docker.service
    systemctl set-default graphical.target
    systemctl start graphical.target
    # Install Chrome https://access.redhat.com/discussions/917293
    curl -o chrome.rpm https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm
    yum -y install redhat-lsb libXScrnSaver
    yum -y localinstall chrome.rpm
    # Install and configure Go
    wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
    tar zxvf go1.8.3.linux-amd64.tar.gz -C /opt
    sed -i 's|^PATH=.*|PATH=$PATH:$HOME/.local/bin:/opt/go/bin:$HOME/bin/|' .bash_profile
    # Install and configure s2i
    wget https://github.com/openshift/source-to-image/releases/download/v1.1.7/source-to-image-v1.1.7-226afa1-linux-amd64.tar.gz
    tar zxvf source-to-image-v1.1.7-226afa1-linux-amd64.tar.gz -C /usr/bin/ ./s2i
    # Fetch Origin CLI
    wget https://github.com/openshift/origin/releases/download/v1.5.1/openshift-origin-client-tools-v1.5.1-7b451fc-linux-64bit.tar.gz
    tar zxvf openshift-origin-client-tools-v1.5.1-7b451fc-linux-64bit.tar.gz --strip-components=1 -C /usr/bin/ openshift-origin-client-tools-v1.5.1-7b451fc-linux-64bit/oc 2>/dev/null
  SHELL
end
