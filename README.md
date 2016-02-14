# Why doing this?
The initiative is to save money and time. I only want a VPN when I am at home.
This means that I need to find a way to easily start and stop EC2 instances on my
own computer. If you do not care the money or use the VPN all day, you do not
need this tool.

# Pre-requisites
* AWS CLI
* Shadowsocks client
* Create a AWS EC2 instance to run shadowsocks server

## How to install and configure AWS CLI
Please refer to the official site.

## How to install Shadowsocks Client
Run the following commands:
```
sudo apt-get install python-pip
pip install shadowsocks
```

## How to create and configure EC2 instance
Please refer to official site on creation of EC2 instance. I created a t2.nano
instance with Ubuntu 14.04 installed. Then I installed the shadowsocks on EC2
instance(same procedures as *Pre-requisites* section). Next I appended command
`sudo ssserver -k "password" -d start` to **/ect/init.d/rc.local**. This line of
code configures the shadowsocks server running when EC2 instance is started.

# Configure the variables
Open main.go with any editor you like. There are 3 variables `instance_id`,
`psd_text` and `local_port`.Change the values of these 3 variables according to
your own configuration. Then build the project as specified at **Build the project**
section.

# Build the project
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
