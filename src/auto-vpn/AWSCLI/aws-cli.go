package AWSCLI

import (
  "os/exec"
  "encoding/json"
  "time"
  "errors"
)

const (
  start_instance = "start-instances"
  stop_instance = "stop-instances"
  describe_instance = "describe-instances"
  no_query = ""
)

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

type AWSEC2Instance struct {
  CachedState DescribeResult //! the latest cached instance state
  InstanceId string
}

func NewEC2Instance(InstId string) *AWSEC2Instance {
  return &AWSEC2Instance{InstanceId: InstId}
}

func (self *AWSEC2Instance) StartInstance() (*StartResult, error) {
  out, err := self.control_ec2(start_instance, no_query)
  if err != nil {
    return nil, err
  }

  var sr StartResult
  err = json.Unmarshal(out, &sr)
  if err != nil {
    return nil, err
  }

  return &sr, nil
}

func (self *AWSEC2Instance) StopInstance() (*StopResult, error) {
  out, err := self.control_ec2(stop_instance, no_query)
  if err != nil {
    return nil, err
  }

  var sr StopResult
  err = json.Unmarshal(out, &sr)
  if err != nil {
    return nil, err
  }

  return &sr, nil
}

func (self *AWSEC2Instance) InstanceState() (*DescribeResult, error) {
  out, err := self.control_ec2(describe_instance, no_query)
  if err != nil {
    return nil, err
  }

  var dr DescribeResult
  err = json.Unmarshal(out, &dr)
  if err != nil {
    return nil, err
  }

  return &dr, nil
}

type StateCallbackT func(error)
func (self *AWSEC2Instance) WaitFor(state string,
  time_out time.Duration, callback StateCallbackT) {
  request_count := int((time_out / time.Second) + 1)
  for i := 0; i < request_count; i++ {
    dr, err := self.InstanceState()
    if err == nil {
      self.CachedState = *dr
      if dr.Reservations[0].Instances[0].State.Name == state {
        callback(nil)
        break;
      }
    }
  }

  callback(errors.New("time out"))
}

func (self *AWSEC2Instance) control_ec2(sub_cmd, query_string string) ([]byte, error) {
  args := []string {
    "ec2", sub_cmd,
    "--instance-id", self.InstanceId,
  }

  if len(query_string) > 0 {
    args = append(args, "--query", query_string)
  }

  cmd := exec.Command("aws", args...)
  out, err := cmd.Output()
  if err != nil {
    return nil, err
  }

  return out, nil
}
