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

	for n := 0; n <= 5; n++ {
		snapshotsOutput, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{})
		if err != nil {
			return "", fmt.Errorf("Failed to retrieve snapshots: %s", err.Error())
		}

		for _, snapshot := range snapshotsOutput.Snapshots {
			if strings.Contains(*snapshot.Description, imageId) {
				return *snapshot.SnapshotId, nil
			}
		}
		fmt.Println("Waiting for snapshots...")
		time.Sleep(time.Millisecond * 5000)
	}

	return "", fmt.Errorf("Unable to find snapshot for image with id %s", imageId)
}
