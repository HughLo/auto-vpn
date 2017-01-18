package main

import (
	"AWSCLI"
	"fmt"
	"log"
	"os"
	"ss"
	"time"
)

//Replace this with your own EC2 Instance ID.
var instanceID = "i-xxxxx"

//Replace this with your shadowsocks server password
var psdText = "password"

//Replace this with your own shadowsocks client port number
var localPort = 10801

func main() {
	subCmd := os.Args[1]

	cliCtrl := AWSCLI.NewEC2Instance(instanceID)

	if cliCtrl == nil {
		log.Fatal("Cannot create AWS CLI instance")
	}

	if subCmd == "start" {
		fmt.Printf("Starting AWS Instance %s\n", cliCtrl.InstanceId)

		sr, err := cliCtrl.StartInstance()
		if err != nil {
			fmt.Println("Start EC2 Instace Error")
			log.Fatal(err)
		}

		fmt.Printf("current status:%s \n", sr.StartingInstances[0].CurrentState.Name)

		cliCtrl.WaitFor("running", time.Second*60, func(err error, dr *AWSCLI.DescribeResult) {
			if err == nil {
				fmt.Printf("The instance %s running successfully \n", cliCtrl.InstanceId)
				ssCtrl := SS.New(dr.Reservations[0].Instances[0].PublicDnsName, "127.0.0.1", psdText, "aes-256-cfb", 8388, localPort, 600)
				if ssCtrl == nil {
					log.Fatal("Cannot create shadowsocks client instance. Do not forget to stop the EC instance.")
				}

				err = ssCtrl.Start()
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
	} else if subCmd == "stop" {
		fmt.Printf("Stopping AWS Instance %s\n", cliCtrl.InstanceId)

		sr, err := cliCtrl.StopInstance()
		if err != nil {
			fmt.Println("Stop EC2 Instace Error. Please try it later or stop it on AWS Offical Site.")
			log.Fatal(err)
		}

		fmt.Printf("current status: %s\n", sr.StoppingInstances[0].CurrentState.Name)

		cliCtrl.WaitFor("stopped", time.Second*60, func(err error, dr *AWSCLI.DescribeResult) {
			if err == nil {
				fmt.Printf("The instance %s stopped successfully \n", cliCtrl.InstanceId)
				ssCtrl := SS.DefaultLocal()
				if ssCtrl == nil {
					log.Fatal("Cannot create shadowsocks client instance.")
				}
				err = ssCtrl.Stop()
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Println("Stop shadowsocks client successfully")
				}
			} else {
				log.Println(err)
			}
		})
	} else if subCmd == "status" {
		dr, err := cliCtrl.InstanceState()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("currnet status: %s\n", dr.Reservations[0].Instances[0].State.Name)
		fmt.Printf("public dns name: %s\n", dr.Reservations[0].Instances[0].PublicDnsName)
	} else if subCmd == "stop-dae" {
		ssCtrl := SS.DefaultLocal()
		if ssCtrl == nil {
			log.Fatal("Cannot create shadowsocks client instance")
		}
		err := ssCtrl.Stop()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Stop shadowsocks client successfully")
		}
	}
}
