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
)

func New(ctx context.Context) detached.Detachable {
	return &Aws{
		session: getSession(ctx),
	}
}

type Aws struct {
	session *session.Session
	config  *config
}

type config struct {
	region string
	zone   string
}

func getRegion() string {
	var regions []string
	defaultPartitions := endpoints.DefaultPartitions()
	for _, partition := range defaultPartitions {
		for _, region := range partition.Regions() {
			regions = append(regions, region.ID())
		}
	}
	fmt.Println(regions)
	return regions[0]
}

func getSession(ctx context.Context) session.Session {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(getRegion())})

	zones, err := svc.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to retrieve availability zones: %s", err))
	}
	var zone string
	for _, avZone := range zones.AvailabilityZones {
		if *avZone.State == "available" {
			zone = *avZone.ZoneName
		}
	}
	if zone == "" {
		return errors.New("No zone found")
	}
	return nil
}
