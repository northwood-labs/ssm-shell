package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/northwood-labs/awsutils"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

type (
	// Ec2Instance represents a list of EC2 instances by name tag and instance ID.
	Ec2Instance struct {
		ID   string
		Name string
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
				ID:   *instance.InstanceId,
				Name: *name,
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
