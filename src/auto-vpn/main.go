package main

import (
  "os"
  "fmt"
  "log"
  "time"
  "AWSCLI"
  "ss"
)

//Replace this with your own EC2 Instance ID.
var instance_id string = "i-ff6e595a"
//Replace this with your shadowsocks server password
var psd_text string = "password"
//Replace this with your own shadowsocks client port number
var local_port int = 10801

func main() {
  sub_cmd := os.Args[1]

  cli_ctrl := AWSCLI.NewEC2Instance(instance_id)

  if cli_ctrl == nil {
    log.Fatal("Cannot create AWS CLI instance")
  }

  if sub_cmd == "start" {
    fmt.Printf("Starting AWS Instance %s\n", cli_ctrl.InstanceId)

    sr, err := cli_ctrl.StartInstance()
    if err != nil {
      fmt.Println("Start EC2 Instace Error")
      log.Fatal(err)
    }

    fmt.Printf("current status:%s \n", sr.StartingInstances[0].CurrentState.Name)

    cli_ctrl.WaitFor("running", time.Second * 60, func (err error, dr *AWSCLI.DescribeResult) {
      if err == nil {
        fmt.Printf("The instance %s running successfully \n", cli_ctrl.InstanceId)
        ss_ctrl := SS.DefaultSSLocal()
        if ss_ctrl == nil {
          log.Fatal("Cannot create shadowsocks client instance. Do not forget to stop the EC instance.")
        }
        ss_ctrl.Password = psd_text
        ss_ctrl.ServerAddr = dr.Reservations[0].Instances[0].PublicDnsName
        ss_ctrl.LocalPort = local_port
        err = ss_ctrl.StartDaemon()
        if err != nil {
          fmt.Println("Cannot start shadowsocks client. Do not forget to stop EC2 instance.")
          log.Fatal(err)
        } else {
          fmt.Println("Start shadowsocks client successfully")
        }
      } else {
        fmt.Println(err)
      }
    })
  } else if sub_cmd == "stop" {
    fmt.Printf("Stopping AWS Instance %s\n", cli_ctrl.InstanceId)

    sr, err := cli_ctrl.StopInstance()
    if err != nil {
      fmt.Println("Stop EC2 Instace Error. Please try it later or stop it on AWS Offical Site.")
      log.Fatal(err)
    }

    fmt.Printf("current status: %s\n", sr.StoppingInstances[0].CurrentState.Name)

    cli_ctrl.WaitFor("stopped", time.Second * 60, func (err error, dr *AWSCLI.DescribeResult) {
      if err == nil {
        fmt.Printf("The instance %s stopped successfully \n", cli_ctrl.InstanceId)
        ss_ctrl := SS.DefaultSSLocal()
        if ss_ctrl == nil {
          log.Fatal("Cannot create shadowsocks client instance.")
        }
        err = ss_ctrl.StopDaemon()
        if err != nil {
          log.Fatal(err)
        } else {
          fmt.Println("Stop shadowsocks client successfully")
        }
      } else {
        log.Println(err)
      }
    })
  } else if sub_cmd == "status" {
    dr, err := cli_ctrl.InstanceState()
    if err != nil {
      log.Fatal(err)
    }

    fmt.Printf("currnet status: %s\n", dr.Reservations[0].Instances[0].State.Name)
  } else if sub_cmd == "stop-dae" {
    ss_ctrl := SS.DefaultSSLocal()
    if ss_ctrl == nil {
      log.Fatal("Cannot create shadowsocks client instance")
    }
    err := ss_ctrl.StopDaemon()
    if err != nil {
      log.Fatal(err)
    } else {
      fmt.Println("Stop shadowsocks client successfully")
    }
  }
}
