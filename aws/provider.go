package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/config"
)

const (
	SECURITY_GROUP_CONFIG_FILE = "security_group.json"
)

func New(ctx context.Context) (detached.Detachable, error) {
	return load(ctx)
}

type AWS struct {
	ID            string `json:"id"`
	Provider      string `json:"provider"`
	Region        string `json:"region"`
	Zone          string `json:"zone"`
	SourceImageId string `json:"sourceImageId"`
	InstanceType  string `json:"instanceType"`
	Username      string `json:"username"`

	ImageId         string `json:"imageId"`
	StackID         string `json:"stackId"`
	SecurityGroupID string `json:"securityGroupId"`
	VolumeID        string `json:"volumeId"`
	SnapshotID      string `json:"snapshotId"`

	_ec2 *ec2.EC2                       `json:"-"`
	_cf  *cloudformation.CloudFormation `json:"-"`
}

func load(ctx context.Context) (*AWS, error) {
	var instance AWS
	if config.Exists() {
		err := config.Load(&instance)
		if err != nil {
			return nil, err
		}
		return &instance, nil
	}
	return nil, errors.New("No config found")
}

func (provider *AWS) ec2() *ec2.EC2 {
	if provider.ec2 == nil {
		if provider.Region == "" {
			log.Fatal("Region is not set yet, please configure your environment first. (detached config)")
		}
		provider._ec2 = ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	}
	return provider._ec2
}

func (provider *AWS) cf() *cloudformation.CloudFormation {
	if provider.ec2 == nil {
		if provider.Region == "" {
			log.Fatal("Region is not set yet, please configure your environment first. (detached config)")
		}
		provider._cf = cloudformation.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})
	}
	return provider._cf
}
