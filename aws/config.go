package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func (provider *Aws) Config(ctx context.Context) error {
	return nil

}

func buildInstance(ctx context.Context) (*Aws, error) {
	region := getRegion()
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	zone, err := getZone(ctx, svc)
	if err != nil {
		return nil, err
	}
	return &Aws{
		Provider:      "aws",
		Region:        region,
		Zone:          zone,
		SourceImageId: "ami-7528ab1a",
		InstanceType:  "t2.micro",
		Username:      "kusowski",
	}, nil
}

func getRegion() string {
	var regions []string
	defaultPartitions := endpoints.DefaultPartitions()
	for _, partition := range defaultPartitions {
		for _, region := range partition.Regions() {
			regions = append(regions, region.ID())
		}
	}
	var region string
	prompt := &survey.Select{
		Message: "Choose a region:",
		Options: regions,
	}
	survey.AskOne(prompt, &region, nil)
	return region
}

func getZone(ctx context.Context, svc *ec2.EC2) (string, error) {
	zones, err := svc.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to retrieve availability zones: %s", err))
	}
	var zoneNames []string
	for _, avZone := range zones.AvailabilityZones {
		if *avZone.State == "available" {
			zoneNames = append(zoneNames, *avZone.ZoneName)
		}
	}
	var zone string
	prompt := &survey.Select{
		Message: "Choose a zone:",
		Options: zoneNames,
	}
	survey.AskOne(prompt, &zone, nil)
	if zone == "" {
		return "", errors.New("No zone selected")
	}
	return zone, nil
}
