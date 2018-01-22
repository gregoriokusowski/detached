package aws

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached/config"
)

var SecurityGroupNotFound = errors.New("Security group not found")

const (
	SECURITY_GROUP_CONFIG_FILE = "security_group.json"
)

func (provider *AWS) GetSecurityGroupId(ctx context.Context) (string, error) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	describeSecurityGroupsOutput, err := svc.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{aws.String("detached-security-group")},
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

func (provider *AWS) UpsertSecurityGroup(ctx context.Context) error {
	csvc := cloudformation.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})

	template, err := ioutil.ReadFile(securityGroupTemplateBodyPath())
	if err != nil {
		return err
	}

	_, err = provider.GetSecurityGroupId(ctx)

	if err == SecurityGroupNotFound {
		fmt.Println("Creating detached security group")
		_, err := csvc.CreateStackWithContext(ctx, &cloudformation.CreateStackInput{
			StackName:    aws.String("detached-security-group"),
			TemplateBody: aws.String(string(template)),
		})
		if err != nil {
			return fmt.Errorf("Failed to create security group: %s", err.Error())
		}
	} else {
		fmt.Println("Updating detached security group")
		_, err := csvc.UpdateStackWithContext(ctx, &cloudformation.UpdateStackInput{
			StackName:    aws.String("detached-security-group"),
			TemplateBody: aws.String(string(template)),
		})
		if err != nil {
			if strings.ContainsAny(err.Error(), "No updates are to be performed") {
				return nil
			}
			return fmt.Errorf("Failed to update security group: %s", err.Error())
		}
	}

	return nil
}

func securityGroupTemplateBodyPath() string {
	return filepath.Join(config.AbsConfigFolder(), SECURITY_GROUP_CONFIG_FILE)
}
