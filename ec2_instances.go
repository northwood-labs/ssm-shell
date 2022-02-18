package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/caarlos0/spin"
	"github.com/northwood-labs/awsutils"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

// Ec2Instance represents a list of EC2 instances by name tag and instance ID.
type Ec2Instance struct {
	ID   string
	Name string
}

// Tag represents a list of EC2 instance tags that we want to filter by.
type Tag struct {
	Name       string
	Equals     string
	Contains   string
	StartsWith string
}

// Filter represents a list of EC2 instance filters that we want to apply.
type Filter struct {
	Name   string
	Equals string
}

func getEc2Instances(tags []Tag, filters []Filter) ([]Ec2Instance, error) {
	s := spin.New("Fetching instances %s ")
	s.Set(spin.Box2)
	s.Start()

	defer s.Stop()

	ctx := context.Background()
	retries := 5
	verbose := false

	config, err := awsutils.GetAWSConfig(ctx, *awsRegion, *awsProfile, retries, verbose)
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

	// Apply user filters
	for i := range filters {
		filter := filters[i]

		ffs = append(ffs, types.Filter{
			Name: aws.String(filter.Name),
			Values: func() []string {
				out := []string{}
				parts := strings.Split(filter.Equals, ",")

				for i := range parts {
					part := parts[i]
					out = append(out, *aws.String(part))
				}

				return out
			}(),
		})
	}

	// Apply user tags
	allTags := getTagEquals(tags)

	for i := range allTags {
		tag := allTags[i]

		ffs = append(ffs, types.Filter{
			Name: aws.String("tag:" + tag.Name),
			Values: []string{
				*aws.String(tag.Equals),
			},
		})
	}

	response, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: ffs,
	})
	if err != nil {
		return []Ec2Instance{}, fmt.Errorf("error looking up instances from EC2 API: %w", err)
	}

	// Everything after this point is client-side filtering.
	allContains := getTagContains(tags)
	allStartsWith := getTagStartsWith(tags)

	for r := range response.Reservations {
		reservation := &response.Reservations[r]
		instances := reservation.Instances

		// Super inefficient. I wanna say O(n²)...?
		if len(allContains) > 0 {
			instances = filterInstances(instances, func(instance types.Instance) bool {
				for i := range instance.Tags {
					t := instance.Tags[i]

					for j := range allContains {
						c := allContains[j]

						if *t.Key == c.Name {
							return strings.Contains(*t.Value, c.Contains)
						}
					}
				}

				return false
			})
		}

		// Super inefficient. I wanna say O(n²)...?
		if len(allStartsWith) > 0 {
			instances = filterInstances(instances, func(instance types.Instance) bool {
				for i := range instance.Tags {
					t := instance.Tags[i]

					for j := range allStartsWith {
						c := allStartsWith[j]

						if *t.Key == c.Name {
							return strings.HasPrefix(*t.Value, c.StartsWith)
						}
					}
				}

				return false
			})
		}

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

func filterTags(vs []Tag, f func(Tag) bool) []Tag {
	vsf := make([]Tag, 0)

	for i := range vs {
		v := vs[i]

		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}

func filterInstances(vs []types.Instance, f func(types.Instance) bool) []types.Instance {
	vsf := make([]types.Instance, 0)

	for i := range vs {
		v := vs[i]

		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}

func getTagEquals(tags []Tag) []Tag {
	return filterTags(tags, func(t Tag) bool {
		return t.Equals != ""
	})
}

func getTagContains(tags []Tag) []Tag {
	return filterTags(tags, func(t Tag) bool {
		return t.Contains != ""
	})
}

func getTagStartsWith(tags []Tag) []Tag {
	return filterTags(tags, func(t Tag) bool {
		return t.StartsWith != ""
	})
}
