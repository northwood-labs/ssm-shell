package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/gookit/color"
	cli "github.com/jawher/mow.cli"
	"github.com/northwood-labs/golang-utils/exiterrorf"
	logrus "github.com/sirupsen/logrus"
)

const (
	fileWritable = 0o666
)

var (
	// Make referencable throughout.
	app        *cli.Cli
	err        error
	awsProfile *string
	awsRegion  *string
	instances  []Ec2Instance

	// Color text.
	colorHeader = color.New(color.FgWhite, color.BgBlue, color.OpBold)

	// Logger.
	logger = logrus.New()
	ff     *os.File

	// Buildtime variables.
	commit  string
	date    string
	version string
)

func main() {
	app = cli.App("ssm-shell", `Simplifies opening a shell session on your EC2 instances using AWS Session
Manager. Supports standard AWS environment variables for authentication.
https://go.aws/3LCabH9

See also:

  * https://github.com/99designs/aws-vault
  * https://github.com/fiveai/aws-okta`)

	app.Version("version", fmt.Sprintf(
		"AWS SSM Shell %s (%s_%s)",
		version,
		runtime.GOOS,
		runtime.GOARCH,
	))

	_, ssmShellLog := os.LookupEnv("SSMSHELL_LOG")

	logger.Level = logrus.DebugLevel
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		FullTimestamp: true,
	})

	// You could set this to any `io.Writer` such as a file
	ff, err = os.OpenFile("/tmp/ssm-shell.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileWritable) // lint:allow_666
	if err == nil && ssmShellLog {
		logger.Out = ff
	} else {
		logger.Out = ioutil.Discard
	}

	awsProfile = app.StringOpt("p profile", "", "The AWS profile entry from your AWS CLI configuration.")
	awsRegion = app.StringOpt("r region", "", "The AWS region to which to communicate.")

	app.Command("connect", "Fetch a list of instances to select from.", cmdConnect)
	app.Command("version", "Verbose information about the build.", cmdVersion)

	err = app.Run(os.Args)
	if err != nil {
		exiterrorf.ExitErrorf(err)
	}
}
