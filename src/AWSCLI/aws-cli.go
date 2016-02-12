//Package AWSCLI wraps the AWS Commnad Line tool usage. To provide a simple
//access to the AWS EC2 Instance. Currently package AWSCLI only support start,
//stop and describe commands.
package AWSCLI

import (
  "os/exec"
  "encoding/json"
  "time"
  "errors"
)

//Define the AWS CLI sub-command string
const (
  start_instance = "start-instances"
  stop_instance = "stop-instances"
  describe_instance = "describe-instances"
  no_query = ""
)

//EC2 Instace state
type Status struct {
  Code int //state code
  Name string //state code string
}

//Instance state returned by "start-instance" and "stop-instance" sub-command
type StartInstanceStatus struct {
  InstanceId string
  CurrentState Status
  PreviousState Status
}

//Return type of StartInstance function
type StartResult struct {
  StartingInstances []StartInstanceStatus
}

//Return type of StopInstance function
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

//Return type of InstaceState function
type DescribeResult struct {
  Reservations []ReservationStruct
}

//An AWSEC2Instance represents an AWS EC2 Instance.
type AWSEC2Instance struct {
  CachedState DescribeResult //! the latest cached instance state
  InstanceId string //! Instance ID
}

//NewEC2Instance creates a new AWSEC2Instance. The InstId variable sets the EC2
//Instace ID. The EC2 Instance ID can be obtained from your AWS Management Console.
func NewEC2Instance(InstId string) *AWSEC2Instance {
  return &AWSEC2Instance{InstanceId: InstId}
}

//StartInstance starts the EC Instance. If error happens, *StartResult will be
//nil. There is no guarantee that the EC2 Instance will be actually running when
//StartInstance returns. Function WaitFor can be used to block execution until
//the state transition finished.
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

//StopInstance stops the EC Instance. If error happens, *StopResult will be
//nil. There is no guarantee that the EC2 Instance will be actually stopped when
//StopInstance returns. Function WaitFor can be used to block execution until
//the state transition finished.
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

//InstanceState returns the current EC2 Instance state. *DescribeResult will be
//nil if error happens.
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

//StateCallbackT represents the callback function type called while the state
//transition inished or timed out.
type StateCallbackT func(error, *DescribeResult)

//WaitFor blocks the execution until the state transition finished or timed out.
//The argument state specifies the desired state. The argument time_out specifies
//the time out limitation. The argument callback will be called either the EC2
//instance transit to the desired state or the time out limitation is exceeded.
func (self *AWSEC2Instance) WaitFor(state string,
  time_out time.Duration, callback StateCallbackT) {
  request_count := int((time_out / time.Second) + 1)
  for i := 0; i < request_count; i++ {
    dr, err := self.InstanceState()
    if err == nil {
      self.CachedState = *dr
      if dr.Reservations[0].Instances[0].State.Name == state {
        callback(nil, dr)
        return
      }
    }
  }

  callback(errors.New("time out"), nil)
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
