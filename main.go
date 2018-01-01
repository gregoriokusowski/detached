package main

import "fmt"

var longDesc = `detached is a tool to create, manage and use remote development environments.

Usage:

        detached command [arguments]

The commands are:

        bootstrap   initialize your configuration
        status      check your current configuration and remote setup
        attach      attach a session to your remote machine, creating it if needed

Use "detached help [command]" for more information about a command.`

func main() {
	fmt.Println(longDesc)
}
