package main

import (
	"fmt"
	"os"

	"github.com/gregoriokusowski/detached"
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
	detached.Bootstrap()
	fmt.Println(os.Args)
	fmt.Println(longDesc)
}
