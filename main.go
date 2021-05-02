package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	mySession := session.Must(session.NewSession())
	svc := ecs.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	subnet := os.Getenv("SUBNET")
	securityGroup := os.Getenv("SECURITY_GROUP")
	taskNetworkConfiguration := ecs.NetworkConfiguration{
		AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
			Subnets:        []*string{aws.String(subnet)},
			SecurityGroups: []*string{aws.String(securityGroup)},
			AssignPublicIp: aws.String("ENABLED"),
		},
	}

	containerName := os.Getenv("CONTAINER_NAME")
	tfPath := os.Getenv("TF_PATH")
	tfCommand := os.Getenv("TF_COMMAND")
	taskContainerOverrides := ecs.ContainerOverride{
		Name: aws.String(containerName),
		Environment: []*ecs.KeyValuePair{
			&ecs.KeyValuePair{
				Name:  aws.String("TF_PATH"),
				Value: aws.String(tfPath),
			},
			&ecs.KeyValuePair{
				Name:  aws.String("TF_COMMAND"),
				Value: aws.String(tfCommand),
			}},
	}

	ecsCluster := os.Getenv("ECS_CLUSTER")
	taskDefinition := os.Getenv("TASK_DEFINITION")
	taskRevision := os.Getenv("TASK_REVISION")
	res, err := svc.RunTask(&ecs.RunTaskInput{
		Cluster:              aws.String(ecsCluster),
		LaunchType:           aws.String("FARGATE"),
		TaskDefinition:       aws.String(fmt.Sprintf("%s:%s", taskDefinition, taskRevision)),
		NetworkConfiguration: &taskNetworkConfiguration,
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{&taskContainerOverrides},
		},
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
