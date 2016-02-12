# auto-vpn
Go to the root folder of your codes and run the following commands to build the executable.
```
export GOPATH=$PWD
go install auto-vpn
```
a **bin** folder will be generated under your **$GOPATH**. The following commands
are provided:
* auto-vpn start
* auto-vpn stop
* auto-vpn status

## auto-vpn start
start AWS instance and sslocal

## auto-vpn stop
stop AWS instance and sslocal

## auto-vpn status
print the current AWS status and its publie DNS name if it exists.

# Pre-requisites
* AWS CLI
* Shadowsocks client
* Create a AWS EC2 instance to run shadowsocks server

## How install and configure AWS CLI
Please refer to the official site.

## How to install Shadowsocks Client
Run the following commands:
```
sudo apt-get install python-pip
pip install shadowsocks
```

## How to create and configure EC2 instance
Please refer to official site on creation of EC2 instance. I created a t2.nano
instance with Ubuntu 14.04 installed. Then I modified **/ect/init.d/rc.local** to
add `sudo ssserver -k "password" -d start` at the end. This line of code configures
the shadowsocks server running when EC2 instance is started.

# How to use
