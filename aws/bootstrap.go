package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	USERDATA = `
#! /sh...
sudo adduser USERNAME --disabled-password
sudo su - USERNAME
echo 'BASE64PUBLICKEY' > ~/.ssh/authorized_keys
`
)

type Config struct {
}

func (provider *Aws) Bootstrap(ctx context.Context) error {

	securityGroupId := ""

	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
	// maybe create current user 	user, _ := user.Current()
	initScript := ""

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	reservation, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          imageId,
		InstanceType:     instanceType,
		MaxCount:         1,
		MinCount:         1,
		SecurityGroupIds: []*string{securityGroupId},
		UserData:         initScript,
	})
	if err != nil {
		return fmt.Errorf("Failed to launch EC2 instance: %s", err.Error())
	}

	instanceId := reservation.Instances[0].InstanceId
	publicIp := reservation.Instances[0].PublicIpAddress
	ebs := reservation.Instances[0].BlockDeviceMappings[0].Ebs
	ebs.SetDeleteOnTermination(false)

	var size int64 = 10
	volumeType := "zgp2"

	// input := &ec2.CreateVolumeInput{
	// 	// AvailabilityZone: aws.String(zone),
	// 	Size:       aws.Int64(size),
	// 	VolumeType: aws.String(volumeType),
	// 	Encrypted:  aws.Bool(true),
	// 	TagSpecifications: []*ec2.TagSpecification{
	// 		&ec2.TagSpecification{
	// 			ResourceType: aws.String("volume"),
	// 			Tags: []*ec2.Tag{
	// 				&ec2.Tag{
	// 					Key:   aws.String("source"),
	// 					Value: aws.String("detached"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// result, err := svc.CreateVolume(input)
	// if err != nil {
	// 	return err
	// }

	// volumeId := *result.VolumeId

	// fmt.Println(volumeId)
	// fmt.Println(result)
	return nil
}
