package aws

import (
	"context"
	"os"
	"os/exec"
)

// func (provider *Aws) SpinUp(ctx context.Context) {
// 	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(provider.Region)})

// 	svc.RequestSpotInstancesWithContext(ctx, &ec2.RequestSpotInstancesInput{

// 	})
// }
func (provider *AWS) Attach(context context.Context) error {
	cmd := exec.Command("tmux", "a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
	return nil
}
