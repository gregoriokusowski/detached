package aws

import (
	"context"
	"errors"

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
