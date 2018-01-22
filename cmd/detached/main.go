package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"github.com/gregoriokusowski/detached"
	"github.com/gregoriokusowski/detached/aws"
)

var longDesc = `detached is a tool to create, manage and use remote development environments.

Usage:

        detached command [arguments]

The commands are:

        config      Creates basic configuration files
        bootstrap   Initialize your configuration
        status      Check your current configuration and remote setup
        attach      Attach a session to your remote machine, creating it if needed

Use "detached help [command]" for more information about a command.`

func main() {
	ctx := context.TODO()
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "config":
			i := &aws.AWS{}
			err := i.Config(ctx)
			if err != nil {
				log.Fatal(err)
			}
		case "bootstrap":
			err := instance(ctx).Bootstrap(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Println(longDesc)
	}
}

func instance(ctx context.Context) detached.Detachable {
	i, err := aws.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
