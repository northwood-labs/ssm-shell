package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/northwood-labs/awsutils"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

type (
	// Ec2Instance represents a list of EC2 instances by name tag and instance ID.
	Ec2Instance struct {
		LaunchTime   time.Time
		ID           string
		Name         string
		Architecture string
		Hypervisor   string
		ImageID      string
		InstanceType string
		Platform     string
		EbsOptimized bool
		EnaSupport   bool
	}

	// Tag represents a list of EC2 instance tags that we want to filter by.
	Tag struct {
		Name       string
		Equals     string
		Contains   string
		StartsWith string
	}

	// Filter represents a list of EC2 instance filters that we want to apply.
	Filter struct {
		Name   string
		Equals string
	}
)

func GetEC2Instances() ([]Ec2Instance, error) {
	ctx := context.Background()
	retries := 5
	verbose := false

	config, err := awsutils.GetAWSConfig(ctx, "", "", retries, verbose)
	if err != nil {
		exiterrorf.ExitErrorf(err)
	}

	var collectedInstances []Ec2Instance

	ec2Client := ec2.NewFromConfig(config)

	// Base filter
	ffs := []types.Filter{
		{
			// Only running instances...
			Name: aws.String("instance-state-name"),
			Values: []string{
				*aws.String("running"),
			},
		},
	}

	response, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: ffs,
	})
	if err != nil {
		return []Ec2Instance{}, fmt.Errorf("error looking up instances from EC2 API: %w", err)
	}

	for r := range response.Reservations {
		reservation := &response.Reservations[r]
		instances := reservation.Instances

		for i := range instances {
			instance := &instances[i]

			// If the conditions exist, apply them.
			name := findName(instance)

			collectedInstances = append(collectedInstances, Ec2Instance{
				ID:           *instance.InstanceId,
				Name:         *name,
				Architecture: string(instance.Architecture),
				Hypervisor:   string(instance.Hypervisor),
				ImageID:      *instance.ImageId,
				InstanceType: string(instance.InstanceType),
				Platform: func() string {
					v := string(instance.Platform)
					if v == "" {
						return "linux"
					} else {
						return v
					}
				}(),
				EbsOptimized: *instance.EbsOptimized,
				EnaSupport:   *instance.EnaSupport,
			})
		}
	}

	return collectedInstances, nil
}

// Calling this is duplicate work. Refactor to collect this data in a single pass.
func findName(instance *types.Instance) *string {
	emptyString := ""

	for t := range instance.Tags {
		tag := instance.Tags[t]

		if *tag.Key == "Name" {
			return tag.Value
		}
	}

	return &emptyString
}
