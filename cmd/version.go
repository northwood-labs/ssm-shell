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
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/northwood-labs/golang-utils/archstring"
	"github.com/spf13/cobra"
)

var (
	// Version represents the version of the software.
	Version = "dev"

	// Commit represents the git commit hash of the software.
	Commit = vcs("vcs.revision", "unknown")

	// BuildDate represents the date the software was built.
	BuildDate = vcs("vcs.time", "unknown")

	// Dirty represents whether or not the git repo was dirty when the software was built.
	Dirty = vcs("vcs.modified", "unknown")

	// PGOEnabled represents whether or not the build leveraged Profile-Guided Optimization (PGO).
	PGOEnabled = vcs("-pgo", "false")

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Long-form version information",
		Long: `Long-form version information, including the build commit hash, build date, Go
version, and external dependencies.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(style.Render(" BUILD INFO "))

			t := table.New().
				Border(lipgloss.RoundedBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
				BorderColumn(true).
				StyleFunc(func(row, col int) lipgloss.Style {
					return lipgloss.NewStyle().Padding(0, 1)
				}).
				Headers("FIELD", "VALUE")

			t.Row("Version", Version)
			t.Row("Go version", runtime.Version())
			t.Row("Git commit", Commit)
			if Dirty == "true" {
				t.Row("Dirty repo", Dirty)
			}
			if !strings.Contains(PGOEnabled, "false") {
				t.Row("PGO", filepath.Base(PGOEnabled))
			}
			t.Row("Build date", BuildDate)
			t.Row("OS/Arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
			t.Row("System", archstring.GetFriendlyName(runtime.GOOS, runtime.GOARCH))
			t.Row("CPU cores", fmt.Sprintf("%d", runtime.NumCPU()))

			fmt.Println(t.Render())

			//----------------------------------------------------------------------

			if buildInfo, ok := debug.ReadBuildInfo(); ok {
				fmt.Println(style.Render(" DEPENDENCIES "))

				td := table.New().
					Border(lipgloss.RoundedBorder()).
					BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
					BorderColumn(true).
					StyleFunc(func(row, col int) lipgloss.Style {
						return lipgloss.NewStyle().Padding(0, 1)
					}).
					Headers("DEPENDENCY", "VERSION")

				for i := range buildInfo.Deps {
					dependency := buildInfo.Deps[i]
					td.Row(dependency.Path, dependency.Version)
				}

				fmt.Println(td.Render())
			}

			fmt.Println("")
		},
	}
)

func init() { // lint:allow_init
	rootCmd.AddCommand(versionCmd)
}

func vcs(key, fallback string) string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for i := range info.Settings {
			setting := info.Settings[i]

			if setting.Key == key {
				return setting.Value
			}
		}
	}

	return fallback
}
