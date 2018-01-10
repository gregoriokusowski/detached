package aws

import (
	"context"
)

const (
	USERDATA = `
#! /sh...
sudo adduser USERNAME --disabled-password
sudo su - USERNAME
echo 'BASE64PUBLICKEY' > ~/.ssh/authorized_keys
`
)

func (provider *Aws) Bootstrap(ctx context.Context) error {

	// svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})

	// fmt.Println("Copying image with encryption to create an Encrypted EBS Volume")
	// copyImageOutput, err := svc.CopyImageWithContext(ctx, &ec2.CopyImageInput{
	// 	SourceImageId: provider.ImageId,
	// 	SourceRegion:  provider.Region,
	// 	Encrypted:     true,
	// 	Name:          "detached-ami",
	// })
	// if err != nil {
	// 	return fmt.Errorf("Failed to create encrypted image: %s", err.Error())
	// }

	// generatedImageId := copyImageOutput.ImageId
	// fmt.Printf("Image %s created\n", generatedImageId)

	// // https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
	// // maybe create current user 	user, _ := user.Current()
	// initScript := ""

	// fmt.Println("Creating dummy instance to spin up volume")
	// reservation, err := svc.RunInstances(&ec2.RunInstancesInput{
	// 	ImageId:          copyImageOutput.ImageId,
	// 	InstanceType:     instanceType,
	// 	MaxCount:         1,
	// 	MinCount:         1,
	// 	SecurityGroupIds: []*string{securityGroupId},
	// 	UserData:         initScript,
	// })
	// if err != nil {
	// 	return fmt.Errorf("Failed to launch EC2 instance: %s", err.Error())
	// }

	// instanceId := reservation.Instances[0].InstanceId
	// ebs := reservation.Instances[0].BlockDeviceMappings[0].Ebs

	// fmt.Println("Detaching volume from instance")
	// volumeAttachment, err := svc.DetachVolume(&ec2.DetachVolumeInput{
	// 	InstanceId: instanceId,
	// 	VolumeId:   ebs.VolumeId,
	// })
	// if err != nil {
	// 	return fmt.Errorf("Failed to detach volume: %s", err.Error())
	// }

	// fmt.Println("Terminating dummy instance")
	// _, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
	// 	InstanceIds: []*string{instanceId},
	// })
	// if err != nil {
	// 	return fmt.Errorf("Failed to terminate dummy instance: %s", err.Error())
	// }

	// // publicIp := reservation.Instances[0].PublicIpAddress
	// // ebs.SetDeleteOnTermination(false)

	// fmt.Println("Modifying volume")
	// _, err = svc.ModifyVolume(&ec2.ModifyVolumeInput{
	// 	VolumeId:   ebs.VolumeId,
	// 	VolumeType: "gp2",
	// 	Size:       10,
	// })
	// if err != nil {
	// 	return fmt.Errorf("Failed to modify volume: %s", err.Error())
	// }

	return nil
}
