package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
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
		ClientToken:   aws.String(provider.ID),
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

	var instanceID, publicDNSName string
	for n := 0; n <= 120; n++ {
		describeInstancesOutput, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{spotRequest.InstanceId},
		})
		if err != nil {
			return fmt.Errorf("Failed to describe instance: %s", err.Error())
		}

		if len(describeInstancesOutput.Reservations) == 0 || len(describeInstancesOutput.Reservations[0].Instances) == 0 {
			time.Sleep(time.Millisecond * 5000)
			continue
		}

		instance := describeInstancesOutput.Reservations[0].Instances[0]
		code := *instance.State.Code
		if code == 0 {
			time.Sleep(time.Millisecond * 5000)
			continue
		}
		if code != 16 {
			return fmt.Errorf("Instance is not pending nor available - got %s", *instance.State.Name)
		}
		instanceID = *instance.InstanceId
		publicDNSName = *instance.PublicDnsName
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
	time.Sleep(time.Millisecond * 5000)

	fmt.Print("Trying to connect...")
	userAndHost := fmt.Sprintf("%s@%s", provider.Username, publicDNSName)
	for {
		cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=5", userAndHost, "exit")
		err := cmd.Run()
		if err == nil {
			fmt.Println("\nConnection is ready.")
			break
		}
		time.Sleep(time.Millisecond * 1000)
		fmt.Print(".")
	}

	cmd := exec.Command("ssh", userAndHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	fmt.Printf("Finishing connection to %s@%s\n", provider.Username, publicDNSName)

	return nil
}
