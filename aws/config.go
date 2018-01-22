package aws

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	defaultSshPublicKeyLocation = ".ssh/id_rsa.pub"
)

// Config prompts for user input and generates the base files for AWS.
func (provider *AWS) Config(ctx context.Context) error {
	region := getRegion()

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	zone, err := getZone(ctx, svc)
	if err != nil {
		return err
	}

	username := getUsername()

	sourceImageId, err := getSourceImageId(ctx, svc)
	if err != nil {
		return err
	}

	instanceType, err := getInstanceType(ctx, svc)
	if err != nil {
		return err
	}

	i := &AWS{
		Provider:      "aws",
		Region:        region,
		Zone:          zone,
		SourceImageId: sourceImageId,
		InstanceType:  instanceType,
		Username:      username,
	}

	publicKey, err := getPublicKey()
	if err != nil {
		return err
	}

	bu := BOOTSTRAP_USERDATA
	bu = strings.Replace(bu, "USERNAME", username, -1)
	bu = strings.Replace(bu, "PUBLIC_KEY", publicKey, -1)
	err = config.AddConfig("bootstrap", bu)
	if err != nil {
		return err
	}

	err = config.AddConfig("spot", SPOT_USERDATA)
	if err != nil {
		return err
	}

	err = config.AddConfig("security_group.json", CLOUDFORMATION_SECURITY_GROUP)
	if err != nil {
		return err
	}

	return config.Save(i)
}

// Given the available AWS regions, prompts a selection.
// TODO: Use a service like http://www.cloudping.info/ and display latencies.
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

// Based on the selected region, prompts for a region.
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

// getUsername prompts for the username, with default value of the current username.
func getUsername() string {
	u, err := user.Current()
	var defaultUsername string
	if err == nil {
		defaultUsername = u.Username
	}
	var username string
	prompt := &survey.Input{
		Default: defaultUsername,
		Message: "Choose your username:",
	}
	survey.AskOne(prompt, &username, nil)
	return username
}

// getSourceImageId currently returns Amazon Linux AMI 2017.09.1.20180115 x86_64 HVM GP2
// This is because Detached only supports amazon linux right now
// TODO: Enable image selection.
func getSourceImageId(ctx context.Context, svc *ec2.EC2) (string, error) {
	return "ami-5652ce39", nil
}

// Should Prompt for the desired instance type.
// TODO: Check if we should enable selection/overwrite during `detached attach`
func getInstanceType(ctx context.Context, svc *ec2.EC2) (string, error) {
	return "t2.micro", nil
}

// Gets the user public key based on path
func getPublicKey() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	defaultPath := filepath.Join(usr.HomeDir, defaultSshPublicKeyLocation)
	var path string
	prompt := &survey.Input{
		Default: defaultPath,
		Message: "Choose your public key path:",
	}
	survey.AskOne(prompt, &path, nil)

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to load public key: %s", err)
	}
	return string(raw), nil
}
