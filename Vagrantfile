# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-16.04"

  config.vm.network "private_network", type: "dhcp"

  config.vm.define "webserver1", primary: true do |server|
    server.vm.network "private_network", ip: "192.168.20.2"
    server.vm.provider "virtualbox" do |vb|
      vb.memory = "1024"
      vb.cpus = "1"
    end
    server.vm.hostname = "webserver1"
    server.vm.provision "shell", inline: <<-SHELL
      echo "Hej from server one!" > /var/www/html/index.html
    SHELL
  end

  config.vm.define "webserver2" do |client|
    client.vm.network "private_network", ip: "192.168.20.3"
    client.vm.provider "virtualbox" do |vb|
      vb.memory = "1024"
      vb.cpus = "1"
    end
    client.vm.hostname = "webserver2"
    client.vm.provision "shell", inline: <<-SHELL
      echo "Hej from server 2!" > /var/www/html/index.html
    SHELL
  end

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    sudo apt-get update
    sudo apt-get -y install apache2
  SHELL
end