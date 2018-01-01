package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (provider *Aws) Bootstrap() error {
	s := session.New()
	svc := ec2.New(s)
	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("us-east-1a"),
		Size:             aws.Int64(10),
		VolumeType:       aws.String("gp2"),
		Encrypted:        aws.Bool(false),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("source"),
						Value: aws.String("detached"),
					},
				},
			},
		},
	}

	result, err := svc.CreateVolume(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
