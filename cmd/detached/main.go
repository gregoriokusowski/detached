package main

import (
	"fmt"
	"os"

	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/aws"
)

var longDesc = `detached is a tool to create, manage and use remote development environments.

Usage:

        detached command [arguments]

The commands are:

        bootstrap   initialize your configuration
        status      check your current configuration and remote setup
        attach      attach a session to your remote machine, creating it if needed

Use "detached help [command]" for more information about a command.`

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		i := instance()
		switch command {
		case "bootstrap":
			err := i.Bootstrap()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else {
		fmt.Println(longDesc)
	}
}

func instance() detached.Detachable {
	return aws.New()
}
