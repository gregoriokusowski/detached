package aws

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
)

var SecurityGroupNotFound = errors.New("Security group not found")

func (provider *AWS) GetSecurityGroupId(ctx context.Context) (string, error) {
	describeSecurityGroupsOutput, err := provider.ec2().DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{aws.String(fmt.Sprintf("detached-security-group-%s", provider.ID))},
	})
	if err != nil {
		if strings.ContainsAny(err.Error(), "InvalidGroup.NotFound") {
			return "", SecurityGroupNotFound
		}
		return "", fmt.Errorf("Failed to get security groups: %s", err.Error())
	}

	if len(describeSecurityGroupsOutput.SecurityGroups) > 0 {
		return *describeSecurityGroupsOutput.SecurityGroups[0].GroupId, nil
	}
	return "", SecurityGroupNotFound
}

func (provider *AWS) CreateSecurityGroupStack(ctx context.Context) (string, error) {
	template, err := ioutil.ReadFile(securityGroupTemplateBodyPath())
	if err != nil {
		return "", err
	}

	fmt.Println("Creating detached security group")
	output, err := provider.cf().CreateStackWithContext(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String(fmt.Sprintf("detached-security-group-%s", provider.ID)),
		TemplateBody: aws.String(string(template)),
	})
	if err != nil {
		return "", fmt.Errorf("Failed to create security group: %s", err.Error())
	}

	return *output.StackId, nil
}

func (provider *AWS) UpdateSecurityGroup(ctx context.Context) error {
	template, err := ioutil.ReadFile(securityGroupTemplateBodyPath())
	if err != nil {
		return err
	}

	fmt.Println("Updating detached security group")
	_, err = provider.cf().UpdateStackWithContext(ctx, &cloudformation.UpdateStackInput{
		StackName:    aws.String("detached-security-group"),
		TemplateBody: aws.String(string(template)),
	})
	if err != nil {
		if !strings.ContainsAny(err.Error(), "No updates are to be performed") {
			return fmt.Errorf("Failed to update security group: %s", err.Error())
		}
	}

	return nil
}

func securityGroupTemplateBodyPath() string {
	return filepath.Join(config.AbsConfigFolder(), SECURITY_GROUP_CONFIG_FILE)
}
