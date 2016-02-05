# auto-vpn
running the following commands to build the executable.
```
export GOPATH=$PWD
go install auto-vpn
```
a **bin** folder will be generated under your **$GOPATH**. The following commands
are provided:
* auto-vpn start
* auto-vpn stop
* auto-vpn status
* auto-vpn sslocal

## auto-vpn start
start AWS instance and sslocal

## auto-vpn stop
stop AWS instance and sslocal

## auto-vpn status
print the current AWS status and its publie DNS name if it exists.

## auto-vpn sslocal
query the AWS instance IP address and run sslocal
