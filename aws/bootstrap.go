package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	USERDATA = `
#!/bin/bash
sudo adduser USERNAME --disabled-password
sudo su - USERNAME
echo 'BASE64PUBLICKEY' > ~/.ssh/authorized_keys
`
)

func (provider *Aws) Bootstrap(ctx context.Context) error {
	fmt.Println("Creating Encrypted AMI")
	imageId, err := provider.CreateEncryptedAMI(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("AMI %s created successfully\n", imageId)

	fmt.Println("Fetching generated Snapshot")
	snapshotId, err := provider.GetSnapshotId(ctx, imageId)
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot %s found for the image %s\n", snapshotId, imageId)

	fmt.Println("Creating security group")
	securityGroupId, err := provider.GetSecurityGroupId(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Security group created successfully")

	fmt.Println("Spinning up one instance to create and setup the volume")
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	_, err = svc.RunInstances(&ec2.RunInstancesInput{
		InstanceType:     aws.String("t2.nano"),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
		SecurityGroupIds: []*string{aws.String(securityGroupId)},
		ImageId:          aws.String(imageId),
		// UserData:         initScript,
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			&ec2.BlockDeviceMapping{
				DeviceName: aws.String("/dev/xvda"),
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(false),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to launch EC2 instance: %s", err.Error())
	}

	return nil
}

func (provider *Aws) CreateEncryptedAMI(ctx context.Context) (string, error) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	copyImageOutput, err := svc.CopyImageWithContext(ctx, &ec2.CopyImageInput{
		Name:          aws.String(fmt.Sprintf("%s detached-copy", provider.ImageId)),
		Encrypted:     aws.Bool(true),
		SourceImageId: aws.String(provider.ImageId),
		SourceRegion:  aws.String(provider.Region),
	})
	if err != nil {
		return "", fmt.Errorf("Failed to copy image: %s", err)
	}

	return *copyImageOutput.ImageId, nil
}

func (provider *Aws) GetSnapshotId(ctx context.Context, imageId string) (string, error) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})

	for n := 0; n <= 20; n++ {
		snapshotsOutput, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{})
		if err != nil {
			return "", fmt.Errorf("Failed to retrieve snapshots: %s", err.Error())
		}

		for _, snapshot := range snapshotsOutput.Snapshots {
			if strings.Contains(*snapshot.Description, imageId) {
				fmt.Printf("State: %s", *snapshot.State)
				if "completed" == *snapshot.State {
					return *snapshot.SnapshotId, nil
				}
			}
		}
		fmt.Println("Waiting for snapshots...")
		time.Sleep(time.Millisecond * 5000)
	}

	return "", fmt.Errorf("Unable to find snapshot for image with id %s", imageId)
}
