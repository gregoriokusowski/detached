package aws

import (
	"context"

	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/config"
)

func New(ctx context.Context) (detached.Detachable, error) {
	return load(ctx)
}

type Aws struct {
	Provider      string `json:"provider"`
	Region        string `json:"region"`
	Zone          string `json:"zone"`
	SourceImageId string `json:"sourceImageId`
	InstanceType  string `json:"instanceType"`
	Username      string `json:"username"`

	ImageId         string `json:"imageId`
	StackID         string `json:"stackId"`
	SecurityGroupID string `json:"securityGroupId"`
	VolumeID        string `json:"volumeId"`
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
