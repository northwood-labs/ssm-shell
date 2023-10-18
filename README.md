# “SSM Shell” for AWS SSM

Simplifies the process of connecting to EC2 instances using AWS Session Manager when you have many instances.

[![Go documentation](https://img.shields.io/badge/pkg.go.dev-documentation-informational?style=flat-square&logo=go)](https://pkg.go.dev/github.com/northwood-labs/ssm-shell)

| Information     | Description                                                                                |
|-----------------|--------------------------------------------------------------------------------------------|
| **Owner(s)**    | [@northwood-labs](https://github.com/northwood-labs), [@skyzyx](https://github.com/skyzyx) |
| **Go Versions** | 1.17                                                                                       |

## Why AWS Session Manager?

SSH is old-school, error-prone, and easy to get wrong.

With the ever-shifting cybersecurity landscape, older ciphers and protocols being cracked over time, and the likelihood of losing SSH keys (or someone stealing them), there are newer, better ways of connecting to EC2 instances in the cloud. AWS Session Manager uses the _AWS Systems Manager_ (SSM) agent to allow you to connect to EC2 instances using the AWS CLI instead of SSH. I'm not going to dive into that here, but here are some links if you don't know what this is:

* <https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html>
* <https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-prerequisites.html>
* <https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-prereqs.html>

If you work for a corporation with lots and lots of AWS accounts, your IT/DevOps/Cloud people are probably taking a look at this if they haven't started using it already.

## Why this project?

Honestly, because I'm tired of having to log into the AWS Console to find the EC2 Instance ID before I pass it to the AWS CLI. Secondly, using the web interface in the AWS Console is _OK_, but I prefer to use the right tool for the job — my terminal.

```bash
aws ssm start-session --target i-abcdef123456
```

Given valid AWS credentials, this will hit the EC2 API first to retrieve a list of running instances, then help you select the instance to which to connect.

## Technical Prerequisites

### AWS CLI v2

If you are on macOS, this is as simple as:

```bash
brew install awscli
```

### AWS Session Manager Plugin

This software plugs into the AWS CLI, allowing you to connect to the instances using it. If you are on macOS, this is a 2-step process.

```bash
brew install session-manager-plugin
```

Next, you need to understand that macOS has an agent called [Gatekeeper](https://support.apple.com/en-us/HT202491) which prevents malware by requiring applications to be [notarized](https://developer.apple.com/news/?id=10032019a). The version of the package vended by Homebrew is not notarized. The version downloaded directly from AWS’s website **is**.

(Why hasn't AWS taken distribution in standard OS package managers into their own hands?)

If you prefer to use Homebrew instead of downloading from the AWS website (like me), you will need to adjust the quarantine settings on the plugin.

```bash
sudo xattr -r -d com.apple.quarantine /usr/local/bin/session-manager-plugin
```

### AWS Vault, AWS Okta, or similar

The **AWS CLI** is a command-line tool for interacting with AWS services. Credentials stored by AWS CLI can also be used with third-party tools which are built using the AWS SDKs. However, AWS CLI **sucks** at making those credentials available to tools other than itself.

[**AWS Vault**](https://github.com/99designs/aws-vault) simplifies this process by communicating with AWS SSO (or your `~/.aws/config` file) up-front, so that you can more easily pass credentials to not just the AWS CLI, but also to any third-party tools which understand AWS credentials. When you regularly manage credentials across multiple AWS accounts, AWS Vault becomes a veritiable necessity.

[**AWS Okta**](https://github.com/fiveai/aws-okta) works similarly, but focuses on vending credentials to (human) users who authenticate with AWS via Okta SSO. (It is also dramtically superior to [Nike’s “Gimme AWS Creds”](https://github.com/Nike-Inc/gimme-aws-creds) tool.)

## Install as a CLI tool

1. You must have the Golang toolchain installed first.

    ```bash
    brew install go
    ```

1. Add `$GOPATH/bin` to your `$PATH` environment variable. By default (i.e., without configuration), `$GOPATH` is defined as `$HOME/go`.

    ```bash
    export PATH="$PATH:$GOPATH/bin"
    ```

1. Once you've done everything above, you can use `go install`.

    ```bash
    go install github.com/northwood-labs/ssm-shell@latest
    ```

## Usage

This will fetch the list of running instances over the EC2 API and respond with an interactive prompt. Hit `Tab` to pop into the interactive UI, and scroll until you find the instance you want. Press `Return`, and you will connect to that instance in your terminal.

<img src="https://github.com/northwood-labs/ssm-shell/raw/main/screenshot.png" width="886">

Use `Control+D` to exit your session.

### Using AWS Vault, AWS Okta, or similar

Assuming you have all of the things working as designed — EC2 instances with SSM agents, Session Manager permissions, local tools installed, etc. — this will simplify logging into instances.

```bash
aws-vault exec {profile} -- ssm-shell connect
```

### Using default AWS CLI profile

Or, if you have credentials setup with `aws configure` using the `default` profile, you can rely on that as well.

```bash
ssm-shell connect
```

### Filtering by instance attributes

Attribute-based filters are processed server-side, and can be used to reduce the number of instances you are presented with. You can apply zero or more attribute-filters. The list of valid attributes can be found at <https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/describe-instances.html>.

For attributes, the valid conditions are: `==` (equals).

```bash
ssm-shell connect --filter architecture==arm64
ssm-shell connect --filter architecture==arm64 --filter availability-zone==us-west-2b
```

### Filtering by tags

Tag-based filters are processed server-side for `==` (equals) conditions, while `=~` (contains) and `=^` (starts with) conditions are processed client-side. You can apply zero or more tag-filters. Like filters, they can be used to reduce the number of instances you are presented with.

For tags, the valid conditions are: `==` (equals), `=~` (contains), and `=^` (starts with).

```bash
ssm-shell connect --tag Application==github-actions
ssm-shell connect --tag Application==github-actions --tag Name=^testing
ssm-shell connect --tag Application==github-actions --tag Name=^testing --tag Name=~amzn2
```
