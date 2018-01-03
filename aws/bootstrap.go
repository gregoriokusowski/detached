package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type VolumeConfig struct {
	Size int64
	Type string
	ID   string
}

func (provider *Aws) Bootstrap(ctx context.Context) error {

	// imageName := "amzn-ami-hvm-2017.09.1.20171120-x86_64-ebs"
	// imageID := "ami-7528ab1a"

	var size int64 = 10
	volumeType := "zgp2"

	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(zone),
		Size:             aws.Int64(size),
		VolumeType:       aws.String(volumeType),
		Encrypted:        aws.Bool(true),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String("volume"),
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

	volumeId := *result.VolumeId

	fmt.Println(result)
	return nil
}
