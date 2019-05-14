package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

var terraformOptions *terraform.Options
var keyPair *aws.Ec2Keypair

func init() {
	terraformState := flag.String("state-folder", "", "Folder path containing the terraform.tfstate file")
	privateKeyFile := flag.String("private-key", "", "private key file path")

	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	privateKey := string(privateKeyBytes[:])

	terraformOptions = &terraform.Options{
		TerraformDir: *terraformState,
	}

	terraform.Init(&testing.T{}, terraformOptions)

	keyPair = &aws.Ec2Keypair{
		KeyPair: &ssh.KeyPair{
			PrivateKey: privateKey,
		},
		Name:   terraform.Output(&testing.T{}, terraformOptions, "key_name"),
		Region: terraform.Output(&testing.T{}, terraformOptions, "region"),
	}

}

func TestEc2Ssh(t *testing.T) {

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ec2-user",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 3
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to public host %s", publicInstanceIP)

	// Run a simple echo command on the server
	expectedText := "Hello, World"
	command := fmt.Sprintf("echo -n '%s'", expectedText)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})

	// Run a command on the server that results in an error,
	expectedText = "Hello, World"
	command = fmt.Sprintf("echo -n '%s' && exit 1", expectedText)
	description = fmt.Sprintf("SSH to public host %s with error command", publicInstanceIP)

	// Verify that we can SSH to the Instance, run the command and see the output
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {

		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err == nil {
			return "", fmt.Errorf("expected SSH command to return an error but got none")
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}
