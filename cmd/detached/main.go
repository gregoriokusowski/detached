package main

import (
	"fmt"
	"os"

	"context"

	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/aws"
)

var longDesc = `detached is a tool to create, manage and use remote development environments.

Usage:

        detached command [arguments]

The commands are:

		config Creates basic configuration files
        bootstrap   initialize your configuration
        status      check your current configuration and remote setup
        attach      attach a session to your remote machine, creating it if needed

Use "detached help [command]" for more information about a command.`

func main() {
	ctx := context.TODO()
	fmt.Println(aws.Default().UpsertSecurityGroup(ctx))
	fmt.Println(aws.Default().GetSecurityGroupId(ctx))
	fmt.Println(aws.Default().UpsertSecurityGroup(ctx))
}
func xmain() {
	ctx := context.TODO()
	if len(os.Args) > 1 {
		command := os.Args[1]
		i, err := instance(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
		switch command {
		case "bootstrap":
			err := i.Bootstrap(ctx)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else {
		fmt.Println(longDesc)
	}
}

func instance(ctx context.Context) (detached.Detachable, error) {
	return aws.New(ctx)
}
