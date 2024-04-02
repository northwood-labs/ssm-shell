package aws

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	CondEquals     = "=="
	CondContains   = "=~"
	CondStartsWith = "=^"
)

func RunCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...) // lint:allow_possible_insecure

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
