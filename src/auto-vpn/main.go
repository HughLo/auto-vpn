package main

import (
  "os"
  "os/exec"
  "fmt"
  "log"
  "time"
  "encoding/json"
)

var instance_id string = "i-ff6e595a"

type Status struct {
  Code int
  Name string
}

type StartInstanceStatus struct {
  InstanceId string
  CurrentState Status
  PreviousState Status
}

type StartResult struct {
  StartingInstances []StartInstanceStatus
}

type StopResult struct {
  StoppingInstances []StartInstanceStatus
}

type MonitoringState struct {
  State string
}

type SecurityGroupStruct struct {
  GroupName string
  GroupId string
}

type InstanceStruct struct {
  Monitoring MonitoringState
  PublicDnsName string
  State Status
  EbsOptimized bool
  LaunchTime string
  PrivateIpAddress string
  InstanceId string
  ImageId string
  PrivateDnsName string
  KeyName string
  SecurityGroups []SecurityGroupStruct
  ClientToken string
  SubnetId string
  InstanceType string
  Architecture string
}

type ReservationStruct struct {
  OwnerId string
  ReservationId string
  Groups []string
  Instances []InstanceStruct
}

type DescribeResult struct {
  Reservations []ReservationStruct
}

var global_state DescribeResult

//! query AWS server to obtain the current instace status
func query_ec2(inst_id, query_string string) (DescribeResult, error) {
  args := []string{"ec2", "describe-instances", "--instance-id", inst_id}

  if len(query_string) > 0 {
    args = append(args, "--query", query_string)
  }

  cmd := exec.Command("aws", args...)

  out, err := cmd.Output()
  if err != nil {
    log.Fatal(err)
  }

  var result DescribeResult
  err = json.Unmarshal([]byte(out), &result)
  if err != nil {
    log.Fatal(err)
  }

  return result, err
}

func control_ec2(sub_cmd, inst_id string) string {
  cmd := exec.Command("aws", "ec2", sub_cmd, "--instance-id", inst_id)
  out, err := cmd.Output()
  if err != nil {
    log.Fatal(err)
  }
  return string(out)
}

func start_instance() {

  //start the instance
  result := control_ec2("start-instances", instance_id)
  var sr StartResult
  err := json.Unmarshal([]byte(result), &sr)
  if err != nil {
    log.Fatal(err)
  }

  //print the running status
  fmt.Printf("starting instance: %s\n", sr.StartingInstances[0].InstanceId)
  fmt.Printf("current state: %s\n", sr.StartingInstances[0].CurrentState.Name)
  fmt.Printf("please wait...\n")

  const (
    max_count = 60
  )

  global_state = DescribeResult{}

  //wait at most 10 seconds for the instance running
  x := 0
  for x <= max_count {
    time.Sleep(time.Second * 1)
    dr, err := query_ec2(instance_id, "")
    if err == nil {
      state := dr.Reservations[0].Instances[0].State.Name
      fmt.Printf("current status: %s\n", state)
      if state == "running" {
        global_state = dr
        break
      }
    }
    x++
  }

  //print the result
  if x > max_count { //time out
    fmt.Printf("cannot get current status. try \"status\" command later\n")
  } else {
    fmt.Printf("instance %s successfully running\n", instance_id)
  }
}

func stop_instance() {

  //stop the instance
  result := control_ec2("stop-instances", instance_id)
  var sr StopResult
  err := json.Unmarshal([]byte(result), &sr)
  if err != nil {
    log.Fatal(err)
  }

  //print the stopping status
  fmt.Printf("stopping instance: %s\n", sr.StoppingInstances[0].InstanceId)
  fmt.Printf("current state: %s\n", sr.StoppingInstances[0].CurrentState.Name)
  fmt.Printf("please wait...\n")

  const (
    max_count = 60
  )

  global_state = DescribeResult{}

  //wait at most 20 seconds for the instance stopped
  x := 0
  for x <= max_count {
    time.Sleep(time.Second * 1)
    dr, err := query_ec2(instance_id, "")
    if err == nil {
      state := dr.Reservations[0].Instances[0].State.Name
      if state == "stopped" {
        global_state = dr
        break
      }
    }
    x++
  }

  //print the result
  if x > max_count { //time out
    fmt.Printf("cannot get current status. try \"status\" command later\n")
  } else {
    fmt.Printf("instance %s successfully stopped\n", instance_id)
  }
}

func instance_status() {
  fmt.Printf("query instance %s status...\n", instance_id)

  dr, err := query_ec2(instance_id, "")
  if err != nil {
    log.Fatal(err)
  }

  global_state = dr

  fmt.Printf("current state: %s\n", dr.Reservations[0].Instances[0].State.Name)
  fmt.Printf("public DNS name: %s\n", dr.Reservations[0].Instances[0].PublicDnsName)
}

func start_shadowsocks() {
  if global_state.Reservations[0].Instances[0].State.Name != "running" {
    log.Fatal("current instance state is not \"running\". cannot start shadowsocks local")
  }

  args := []string {
    "sslocal",
    "-s", global_state.Reservations[0].Instances[0].PublicDnsName,
    "-p", "8388",
    "-b", "127.0.0.1",
    "-l", "10801",
    "-k", "password",
    "-m", "aes-256-cfb",
    "-t", "600",
    "-d", "start",
  }

  cmd := exec.Command("sudo", args...)
  _, err := cmd.Output()
  if err != nil {
    log.Fatal(err)
  } else {
    fmt.Println("start shadowsocks local successfully")
  }
}

func stop_shadowsocks() {
  args := []string {
    "sslocal",
    "-s", "",
    "-d", "stop",
  }

  cmd := exec.Command("sudo", args...)
  _, err := cmd.Output()
  if err != nil {
    log.Fatal(err)
  } else {
    fmt.Println("stop shadowsocks local successfully")
  }
}

func main() {
  sub_cmd := os.Args[1]

  if sub_cmd == "start" {
    start_instance()
    start_shadowsocks()
  } else if sub_cmd == "stop" {
    stop_instance()
    stop_shadowsocks()
  } else if sub_cmd == "status" {
    instance_status()
  } else if sub_cmd == "sslocal" {
    instance_status()
    start_shadowsocks()
  }
}
