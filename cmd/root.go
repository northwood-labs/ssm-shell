// Copyright 2023–2024, Northwood Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/northwood-labs/ssm-shell/aws"
	"github.com/spf13/cobra"
)

const height = 20

var (
	instanceID string
	instances  []aws.Ec2Instance

	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "ssm-shell",
	})

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "ssm-shell",
		Short: "Simplifies the process of connecting to EC2 Instances using AWS Session Manager.",
		Long: `--------------------------------------------------------------------------------
ssm-shell

Simplifies the process of connecting to EC2 Instances using AWS Session Manager.

Disabling SSH and leveraging AWS Session Manager to connect to EC2 Instances is
the recommended approach for managing EC2 Instances. This approach is more
secure and does not require the need to manage SSH keys.
--------------------------------------------------------------------------------`,
		Run: func(cmd *cobra.Command, args []string) {
			err := spinner.New().
				Title("Getting EC2 instances for this account...").
				Action(func(instances *[]aws.Ec2Instance) func() {
					return func() {
						insts, e := aws.GetEC2Instances()
						if e != nil {
							logger.Fatal(e)
						}

						// Sort by text
						sort.SliceStable(insts, func(i, j int) bool {
							return strings.ToLower(insts[i].Name) < strings.ToLower(insts[j].Name)
						})

						*instances = insts
					}
				}(&instances)).
				Run()

			if err != nil {
				logger.Fatal(err)
			}

			options := []huh.Option[string]{}

			for i := range instances {
				options = append(
					options,
					huh.NewOption(
						fmt.Sprintf(
							"%s (%s)",
							instances[i].Name,
							instances[i].ID,
						),
						instances[i].ID,
					),
				)
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Connect to which SSM-enabled EC2 instance?").
						Options(options...).
						Value(&instanceID).
						WithHeight(height),
				),
			)

			err = form.Run()
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Printf("Connecting to instance %s...\n", instanceID)
			aws.RunCommand(
				strings.Split(
					fmt.Sprintf("aws ssm start-session --target %s", instanceID),
					" ",
				),
			)
		},
	}

	style = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder())
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Fatal(err)
	}
}
