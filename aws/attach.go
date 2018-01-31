package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
)

func (provider *AWS) Attach(ctx context.Context) error {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})

	spotUserData, err := config.GetConfig("spot")
	if err != nil {
		return err
	}

	fmt.Println("Requesting spot instance")
	reqOutput, err := svc.RequestSpotInstancesWithContext(ctx, &ec2.RequestSpotInstancesInput{
		SpotPrice:     aws.String("0.05"),
		InstanceCount: aws.Int64(1),
		LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
			ImageId:          aws.String(provider.ImageId),
			InstanceType:     aws.String(provider.InstanceType),
			SecurityGroupIds: []*string{aws.String(provider.SecurityGroupID)},
			UserData:         aws.String(base64.StdEncoding.EncodeToString(spotUserData)),
			Placement: &ec2.SpotPlacement{
				AvailabilityZone: aws.String(provider.Zone),
			},
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{
				&ec2.BlockDeviceMapping{
					DeviceName: aws.String("/dev/xvda"),
					Ebs: &ec2.EbsBlockDevice{
						DeleteOnTermination: aws.Bool(true),
						VolumeType:          aws.String("gp2"),
						VolumeSize:          aws.Int64(8),
						SnapshotId:          aws.String(provider.SnapshotID),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to request spot instance: %s", err.Error())
	}

	spotRequestID := *reqOutput.SpotInstanceRequests[0].SpotInstanceRequestId

	fmt.Println("Waiting request to be fulfilled")
	instanceRequestOutput, err := svc.DescribeSpotInstanceRequestsWithContext(ctx, &ec2.DescribeSpotInstanceRequestsInput{
		SpotInstanceRequestIds: []*string{aws.String(spotRequestID)},
	})
	if err != nil {
		return fmt.Errorf("Failed to retrieve spot request status: %s", err.Error())
	}

	spotRequest := instanceRequestOutput.SpotInstanceRequests[0]
	state := *spotRequest.State
	if state != "open" && state != "active" {
		return fmt.Errorf("Spot instance request failed. State was %s", state)
	}

	var instanceID string
	for n := 0; n <= 120; n++ {
		instanceStatusOutput, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
			InstanceIds: []*string{spotRequest.InstanceId},
		})
		if err != nil {
			return fmt.Errorf("Failed to describe instance: %s", err.Error())
		}
		if len(instanceStatusOutput.InstanceStatuses) == 0 {
			time.Sleep(time.Millisecond * 5000)
			continue
		}

		code := *instanceStatusOutput.InstanceStatuses[0].InstanceState.Code
		if code == 0 {
			time.Sleep(time.Millisecond * 5000)
			continue
		}
		if code != 16 {
			return fmt.Errorf("Instance is not pending nor available")
		}
		instanceID = *instanceStatusOutput.InstanceStatuses[0].InstanceId
		break
	}

	_, err = svc.AttachVolumeWithContext(ctx, &ec2.AttachVolumeInput{
		InstanceId: aws.String(instanceID),
		VolumeId:   aws.String(provider.VolumeID),
		Device:     aws.String("/dev/xvdf"),
	})
	if err != nil {
		return fmt.Errorf("Failed to attach volume: %s", err.Error())
	}
	fmt.Println("Volume attached. Waiting for disk relabeling and reboot")

	// cmd := exec.Command("tmux", "a")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin

	// cmd.Run()
	return nil
}
