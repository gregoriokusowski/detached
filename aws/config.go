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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
	uuid "github.com/satori/go.uuid"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	defaultSshPublicKeyLocation = ".ssh/id_rsa.pub"
)

// Config prompts for user input and generates the base files for AWS.
func (provider *AWS) Config(ctx context.Context) error {
	id := uuid.NewV4().String()

	provider.Region = getRegion()

	zone, err := provider.getZone(ctx)
	if err != nil {
		return err
	}

	username := getUsername()

	sourceImageId, err := provider.getSourceImageId(ctx)
	if err != nil {
		return err
	}

	instanceType, err := provider.getInstanceType(ctx)
	if err != nil {
		return err
	}

	i := &AWS{
		ID:            id,
		Provider:      "aws",
		Region:        provider.Region,
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

	cf := CLOUDFORMATION_SECURITY_GROUP
	cf = strings.Replace(cf, "DETACHED_ID", id, -1)
	err = config.AddConfig("security_group.json", cf)
	if err != nil {
		return err
	}

	err = config.Save(i)
	if err != nil {
		return err
	}

	fmt.Println("Config is done. You can check it and run bootstrap when ready!")
	return nil
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
func (provider *AWS) getZone(ctx context.Context) (string, error) {
	zones, err := provider.ec2().DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{})
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
func (provider *AWS) getSourceImageId(ctx context.Context) (string, error) {
	availableImages, err := provider.ec2().DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("name"),
				Values: []*string{aws.String("amzn-ami-hvm-2017.09.1.20180115-x86_64-gp2")},
			},
			&ec2.Filter{
				Name:   aws.String("owner-alias"),
				Values: []*string{aws.String("amazon")},
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(availableImages.Images) == 0 {
		return "", fmt.Errorf("Unable to find image")
	}
	return *availableImages.Images[0].ImageId, nil
}

// Should Prompt for the desired instance type.
// TODO: Check if we should enable selection/overwrite during `detached attach`
func (_ *AWS) getInstanceType(_ context.Context) (string, error) {
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
