package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/config"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func New(ctx context.Context) (detached.Detachable, error) {
	return load(ctx)
}

type Aws struct {
	Provider     string `json:"provider"`
	Region       string `json:"region"`
	Zone         string `json:"zone"`
	ImageId      string `json:"imageId`
	InstanceType string `json:"instanceType"`
	Username     string `json:"username"`
	SshPort      int    `json:"sshPort"`
}

func Default() *detached.Detachable {
	return &Aws{
		Provider:     "aws",
		Region:       "eu-central-1",
		Zone:         "eu-central-1-a",
		ImageId:      "ami-7528ab1a",
		InstanceType: "t2.micro",
		Username:     "kusowski",
		SshPort:      22,
	}
}

func load(ctx context.Context) (*Aws, error) {
	var instance *Aws
	if config.Exists() {
		err := config.Load(*instance)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
	instance, err := buildInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = config.Save(instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func buildInstance(ctx context.Context) (*Aws, error) {
	region := getRegion()
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	zone, err := getZone(ctx, svc)
	if err != nil {
		return nil, err
	}
	return &Aws{
		Provider: "aws",
		Region:   region,
		Zone:     zone,
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
