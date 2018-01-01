package detached

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Bootstrap() {
	svc := ec2.New(session.New())
	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("us-east-1a"),
		Size:             aws.Int64(80),
		VolumeType:       aws.String("gp2"),
	}

	result, err := svc.CreateVolume(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}

func Status() {

}

func Attach() {
	cmd := exec.Command("tmux", "a")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

}
