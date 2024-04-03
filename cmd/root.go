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

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	helpText = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2) // lint:allow_raw_number

	style = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder())

	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "ssm-shell",
	})

	keys = keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "make selection"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
	}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "ssm-shell",
		Short: "Simplifies the process of connecting to EC2 Instances using AWS Session Manager.",
		Long: helpText.Render(`ssm-shell

Simplifies the process of connecting to EC2 Instances using AWS Session Manager.

Disabling SSH and leveraging AWS Session Manager to connect to EC2 Instances is
the recommended approach for managing EC2 Instances. This approach is more
secure and does not require the need to manage SSH keys.`),
		Run: func(cmd *cobra.Command, args []string) {
			err := spinner.New().
				Title("Getting EC2 instances for this account...").
				Type(spinner.Dots).
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

			columns := []table.Column{
				{Title: "Name", Width: 35},     // lint:allow_raw_number
				{Title: "ID", Width: 20},       // lint:allow_raw_number
				{Title: "CPU", Width: 6},       // lint:allow_raw_number
				{Title: "Type", Width: 15},     // lint:allow_raw_number
				{Title: "AMI", Width: 25},      // lint:allow_raw_number
				{Title: "Platform", Width: 10}, // lint:allow_raw_number
			}

			rows := []table.Row{}

			for i := range instances {
				rows = append(
					rows,
					table.Row{
						instances[i].Name,
						instances[i].ID,
						instances[i].Architecture,
						instances[i].InstanceType,
						instances[i].ImageID,
						instances[i].Platform,
					},
				)
			}

			t := table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(height),
			)

			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
			s.Selected = s.Selected.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)
			t.SetStyles(s)

			m := model{
				table: t,
				keys:  keys,
				help:  help.New(),
			}
			if _, err := tea.NewProgram(m).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

			if instanceID != "" {
				fmt.Printf("Connecting to instance %s...\n", instanceID)
				aws.RunCommand(
					strings.Split(
						fmt.Sprintf("aws ssm start-session --target %s", instanceID),
						" ",
					),
				)
			}
		},
	}
)

type (
	model struct {
		help     help.Model
		lastKey  string
		keys     keyMap
		table    table.Model
		quitting bool
	}

	// keyMap defines a set of keybindings. To work for help it must satisfy
	// key.Map. It could also very easily be a map[string]key.Binding.
	keyMap struct {
		Up    key.Binding
		Down  key.Binding
		Help  key.Binding
		Enter key.Binding
		Quit  key.Binding
	}
)

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding { // lint:allow_large_memory // Implementing a model I have no control over.
	return []key.Binding{
		k.Help,
		k.Enter,
		k.Quit,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding { // lint:allow_large_memory // Implementing a model I have no control over.
	return [][]key.Binding{
		{ // first column
			k.Up,
			k.Down,
		},
		{ // second column
			k.Help,
			k.Quit,
		},
		{ // third column
			k.Enter,
		},
	}
}

func (m model) Init() tea.Cmd { // lint:allow_large_memory // Implementing a model I have no control over.
	return nil
}

func (m model) Update( // lint:allow_large_memory // Implementing a model I have no control over.
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Enter):
			m.quitting = true
			instanceID = m.table.SelectedRow()[1]

			return m, tea.Quit
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true

			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string { // lint:allow_large_memory // Implementing a model I have no control over.
	if m.quitting {
		return ""
	}

	helpView := m.help.View(m.keys)

	return baseStyle.Render(m.table.View()) + "\n" + helpView
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Fatal(err)
	}
}
