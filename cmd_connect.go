package main

import (
	"fmt"
	"sort"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	cli "github.com/jawher/mow.cli"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

const (
	condEquals     = "=="
	condContains   = "=~"
	condStartsWith = "=^"
)

func cmdConnect(cmd *cli.Cmd) {
	tags := cmd.StringsOpt("t tag", []string{}, fmt.Sprintf(
		"Tag names and tag values, separated by a condition. Conditions are `%s` (equals),\n"+
			"`%s` (contains), and `%s` (starts with). Flag can be called multiple times.",
		condEquals,
		condContains,
		condStartsWith,
	))

	filters := cmd.StringsOpt("f filter", []string{}, fmt.Sprintf(
		"Filter names and filter values, separated by a `%s` (equals) condition. Flag can\n"+
			"be called multiple times. See https://bit.ly/3JqctHs for list of valid values.",
		condEquals,
	))

	cmd.Action = func() {
		tagStructs := processTagInput(*tags)
		filterStructs := processFilterInput(*filters)

		instances, err = getEc2Instances(tagStructs, filterStructs)
		if err != nil {
			exiterrorf.ExitErrorf(err)
		}

		// Sort by text
		sort.SliceStable(instances, func(i, j int) bool {
			return strings.ToLower(instances[i].Name) < strings.ToLower(instances[j].Name)
		})

		instanceName := prompt.Input(
			"The instance to connect to [press tab]: ",
			func(tags []Tag) func(in prompt.Document) []prompt.Suggest {
				return instanceCompleter
			}(tagStructs),
			CustomOptions()...,
		)

		runCommand(
			strings.Split(
				fmt.Sprintf("aws ssm start-session --target %s", instanceName),
				" ",
			),
		)
	}
}

func processTagInput(tags []string) []Tag {
	out := []Tag{}

	for i := range tags {
		tag := tags[i]

		if strings.Contains(tag, condEquals) {
			result := strings.Split(tag, condEquals)

			out = append(out, Tag{
				Name:   result[0],
				Equals: result[1],
			})
		} else if strings.Contains(tag, condContains) {
			result := strings.Split(tag, condContains)

			out = append(out, Tag{
				Name:     result[0],
				Contains: result[1],
			})
		} else if strings.Contains(tag, condStartsWith) {
			result := strings.Split(tag, condStartsWith)

			out = append(out, Tag{
				Name:       result[0],
				StartsWith: result[1],
			})
		}
	}

	return out
}

func processFilterInput(filters []string) []Filter {
	out := []Filter{}

	for i := range filters {
		filter := filters[i]
		result := strings.Split(filter, condEquals)

		out = append(out, Filter{
			Name:   result[0],
			Equals: result[1],
		})
	}

	return out
}

func instanceCompleter(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for i := range instances {
		instance := instances[i]

		s = append(s, prompt.Suggest{
			Text: func() string {
				if instance.ID != "" {
					return instance.ID
				}

				return ""
			}(),
			Description: func() string {
				if instance.Name != "" {
					return instance.Name
				}

				return ""
			}(),
		})
	}

	return prompt.FilterFuzzy(s, in.GetWordBeforeCursorUntilSeparator("\n"), true)
}
