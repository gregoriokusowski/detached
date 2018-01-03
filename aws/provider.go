package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const CONFIG_FOLDER = "~/.detached"
const CONFIG_PATH = "~/.detached/default.json"

func New(ctx context.Context) (detached.Detachable, error) {
	return load(ctx)
}

type Aws struct {
	Provider string `json:"provider"`
	Region   string `json:"region"`
	Zone     string `json:"zone"`
}

func load(ctx context.Context) (*Aws, error) {
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

func persist(instance *Aws) error {
	bytes, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CONFIG_PATH, bytes, 0644)

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
