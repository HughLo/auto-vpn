//Package SS wraps the life-cycle control of shadowsocks client. This package
//uses sslocal script to control the shadowsocks life-cycle
package SS

import (
  "os/exec"
  "strconv"
)

//SSLocal represents a shadowsocks client.
type SSLocal struct {
  ServerAddr string
  ServerPort int
  LocalAddr string
  LocalPort int
  Password string
  Timeout int
  Method string
}

//DefaultSSLocal creates a new SSLocal with default configuration.
func DefaultSSLocal() *SSLocal {
  return &SSLocal {
    ServerAddr: "",
    ServerPort: 8388,
    LocalAddr: "127.0.0.1",
    LocalPort: 1088,
    Password: "",
    Timeout: 600,
    Method: "aes-256-cfb",
  }
}

//New creates a new SSLocal. srv is server address, local is local address, psd
//is the shadowsocks server password, local_port is the port number of shadowsocks
//client, time_out is the time out setup of shadowsocks client.
func New(srv, local, psd, method string,
  srv_port, local_port, time_out int) *SSLocal {
  return &SSLocal{
    ServerAddr: srv,
    ServerPort: srv_port,
    LocalAddr: local,
    LocalPort: local_port,
    Password: psd,
    Timeout: time_out,
    Method: method,
  }
}

//StartDaemon starts sslocal as daemon
func (self *SSLocal) StartDaemon() error {
  return self.control_sslocal_daemon("start")
}

//StopDaemon stops sslocal daemon
func (self *SSLocal) StopDaemon() error {
  return self.control_sslocal_daemon("stop")
}

//RestartDaemon restarts sslocal daemon
func (self *SSLocal) RestartDaemon() error {
  return self.control_sslocal_daemon("restart")
}

func (self *SSLocal) control_sslocal_daemon(dae string) error {
  args := []string {
    "sslocal",
    "-s", self.ServerAddr,
    "-p", strconv.Itoa(self.ServerPort),
    "-b", self.LocalAddr,
    "-l", strconv.Itoa(self.LocalPort),
    "-k", self.Password,
    "-m", self.Method,
    "-t", strconv.Itoa(self.Timeout),
    "-d", dae,
  }

  cmd := exec.Command("sudo", args...)
  _, err := cmd.Output()
  return err
}
