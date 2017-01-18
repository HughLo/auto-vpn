//Package SS wraps the life-cycle control of shadowsocks client. This package
//uses sslocal script to control the shadowsocks life-cycle
package SS

import (
	"os/exec"
	"runtime"
	"strconv"
)

type controlSslocalDaemon interface {
	Control(string) error
}

//Local represents a shadowsocks client.
type Local struct {
	ServerAddr string
	ServerPort int
	LocalAddr  string
	LocalPort  int
	Password   string
	Timeout    int
	Method     string
}

//LocalControl is the interface to control the shadowsocks client
type LocalControl interface {
	Start() error
	Stop() error
	Restart() error
}

//WinLocal represents a shadowsocks client in windows platform.
type WinLocal struct {
	Local
}

//LinuxLocal represents a shadowsocks client in linux platform.
type LinuxLocal struct {
	Local
}

//Control sslocal in windows platform
func (wlocal *WinLocal) control(dae string) error {
	args := []string{
		"-s", wlocal.ServerAddr,
		"-p", strconv.Itoa(wlocal.ServerPort),
		"-b", wlocal.LocalAddr,
		"-l", strconv.Itoa(wlocal.LocalPort),
		"-k", wlocal.Password,
		"-m", wlocal.Method,
		"-t", strconv.Itoa(wlocal.Timeout),
		"-d", dae,
	}

	cmd := exec.Command("sslocal", args...)
	_, err := cmd.Output()
	return err
}

//Start the shadowsocks client
func (wlocal *WinLocal) Start() error {
	return wlocal.control("start")
}

//Stop the shadowsocks client
func (wlocal *WinLocal) Stop() error {
	return wlocal.control("stop")
}

//Restart the shadowsocks client
func (wlocal *WinLocal) Restart() error {
	return wlocal.control("restart")
}

//Control sslocal in windows platform
func (llocal *LinuxLocal) control(dae string) error {
	args := []string{
		"sslocal",
		"-s", llocal.ServerAddr,
		"-p", strconv.Itoa(llocal.ServerPort),
		"-b", llocal.LocalAddr,
		"-l", strconv.Itoa(llocal.LocalPort),
		"-k", llocal.Password,
		"-m", llocal.Method,
		"-t", strconv.Itoa(llocal.Timeout),
		"-d", dae,
	}

	cmd := exec.Command("sudo", args...)
	_, err := cmd.Output()
	return err
}

//Start the shadowsocks client
func (llocal *LinuxLocal) Start() error {
	return llocal.control("start")
}

//Stop the shadowsocks client
func (llocal *LinuxLocal) Stop() error {
	return llocal.control("stop")
}

//Restart the shadowsocks client
func (llocal *LinuxLocal) Restart() error {
	return llocal.control("restart")
}

//DefaultLocal creates a new Local with default configuration.
func DefaultLocal() LocalControl {
	local := Local{
		ServerAddr: "",
		ServerPort: 8388,
		LocalAddr:  "127.0.0.1",
		LocalPort:  1088,
		Password:   "",
		Timeout:    600,
		Method:     "aes-256-cfb",
	}

	var localControl LocalControl

	switch runtime.GOOS {
	case "windows":
		winLocal := new(WinLocal)
		winLocal.Local = local
		localControl = winLocal

	case "linux":
		linuxLocal := new(LinuxLocal)
		linuxLocal.Local = local
		localControl = linuxLocal
	}

	return localControl
}

//New creates a new SSLocal. srv is server address, local is local address, psd
//is the shadowsocks server password, local_port is the port number of shadowsocks
//client, time_out is the time out setup of shadowsocks client.
func New(srv, localAddr, psd, method string,
	srvPort, localPort, timeOut int) LocalControl {
	local := Local{
		ServerAddr: srv,
		ServerPort: srvPort,
		LocalAddr:  localAddr,
		LocalPort:  localPort,
		Password:   psd,
		Timeout:    timeOut,
		Method:     method,
	}

	var localControl LocalControl

	switch runtime.GOOS {
	case "windows":
		winLocal := new(WinLocal)
		winLocal.Local = local
		localControl = winLocal

	case "linux":
		linuxLocal := new(LinuxLocal)
		linuxLocal.Local = local
		localControl = linuxLocal
	}

	return localControl
}
