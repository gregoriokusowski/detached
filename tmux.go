package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("tmux", "a")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
