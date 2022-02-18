package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"text/tabwriter"

	cli "github.com/jawher/mow.cli"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

func cmdVersion(cmd *cli.Cmd) {
	cmd.Action = func() {
		fmt.Println(colorHeader.Render(" BASIC "))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		fmt.Fprintf(w, " Version:\t%s\t\n", version)
		fmt.Fprintf(w, " Go version:\t%s\t\n", runtime.Version())
		fmt.Fprintf(w, " Git commit:\t%s\t\n", commit)
		fmt.Fprintf(w, " Build date:\t%s\t\n", date)
		fmt.Fprintf(w, " OS/Arch:\t%s/%s\t\n", runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(w, " System:\t%s\t\n", getFriendlyName(runtime.GOOS, runtime.GOARCH))
		fmt.Fprintf(w, " CPU Cores:\t%d\t\n", runtime.NumCPU())

		err := w.Flush()
		if err != nil {
			exiterrorf.ExitErrorf(err)
		}

		fmt.Println("")

		//----------------------------------------------------------------------

		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(colorHeader.Render(" DEPENDENCIES "))

			w = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

			for i := range buildInfo.Deps {
				dependency := buildInfo.Deps[i]
				fmt.Fprintf(w, " %s\t%s\t\n", dependency.Path, dependency.Version)
			}
		}

		err = w.Flush()
		if err != nil {
			exiterrorf.ExitErrorf(err)
		}

		fmt.Println("")
	}
}
