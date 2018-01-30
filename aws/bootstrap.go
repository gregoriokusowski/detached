package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
)

func (provider *AWS) Bootstrap(ctx context.Context) error {
	fmt.Println("Creating security group cloudformation stack")
	stackID, err := provider.CreateSecurityGroupStack(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Security group stack created successfully")

	fmt.Println("Creating Encrypted AMI")
	imageId, err := provider.CreateEncryptedAMI(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("AMI %s created successfully and is available\n", imageId)

	fmt.Println("Retrieving security group")
	securityGroupId, err := provider.GetSecurityGroupId(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Security group retrieved successfully")

	bootstrap, err := config.GetConfig("bootstrap")
	if err != nil {
		return err
	}

	fmt.Println("Spinning up one instance to create and setup the volume")
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	_, err = svc.RunInstances(&ec2.RunInstancesInput{
		InstanceType:     aws.String("t2.nano"),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
		SecurityGroupIds: []*string{aws.String(securityGroupId)},
		ImageId:          aws.String(imageId),
		UserData:         aws.String(base64.StdEncoding.EncodeToString(bootstrap)),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			&ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(false),
				},
			},
		},
		InstanceInitiatedShutdownBehavior: aws.String("terminate"),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("source"),
						Value: aws.String("detached"),
					},
				},
			},
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
	})
	if err != nil {
		return fmt.Errorf("Failed to launch EC2 instance: %s", err.Error())
	}

	fmt.Println("Retrieving volume ID")
	var volumeID string
	for n := 0; n <= 120; n++ {
		volumeOutput, err := svc.DescribeVolumesWithContext(ctx, &ec2.DescribeVolumesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("tag:source"),
					Values: []*string{aws.String("detached")},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("Failed to retrieve volume information: %s", err.Error())
		}
		if len(volumeOutput.Volumes) > 0 {
			volumeID = *volumeOutput.Volumes[0].VolumeId
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	if volumeID == "" {
		return fmt.Errorf("Volume ID was not found")
	}

	fmt.Println("Saving config.")
	provider.ImageId = imageId
	provider.SecurityGroupID = securityGroupId
	provider.StackID = stackID
	provider.VolumeID = volumeID
	err = config.Save(provider)

	return err
}

func (provider *AWS) CreateEncryptedAMI(ctx context.Context) (string, error) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	copyImageOutput, err := svc.CopyImageWithContext(ctx, &ec2.CopyImageInput{
		Name:          aws.String(fmt.Sprintf("%s detached-copy", provider.SourceImageId)),
		Encrypted:     aws.Bool(true),
		SourceImageId: aws.String(provider.SourceImageId),
		SourceRegion:  aws.String(provider.Region),
	})
	if err != nil {
		return "", fmt.Errorf("Failed to copy image: %s", err)
	}

	fmt.Print("Waiting for image (it may take a few minutes) ...")
	for n := 0; n <= 120; n++ {
		images, err := svc.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
			ImageIds: []*string{copyImageOutput.ImageId},
		})
		if err != nil {
			return "", fmt.Errorf("Failed to retrieve images: %s", err.Error())
		}

		for _, image := range images.Images {
			fmt.Print(".")
			if "available" == *image.State {
				return *copyImageOutput.ImageId, nil
			}
		}
		time.Sleep(time.Millisecond * 5000)
	}

	return "", errors.New("Image was not available")
}
