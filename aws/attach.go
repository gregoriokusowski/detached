package aws

import (
	"os"
	"os/exec"
)

func (provider *Aws) Attach() error {
	cmd := exec.Command("tmux", "a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
	return nil
}
