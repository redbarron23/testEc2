package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Connection to target
type Connection struct {
	*ssh.Client
}

// Connect to Ec2 Instance
func Connect(addr, user string) (*Connection, error) {

	pemBytes, err := ioutil.ReadFile(pemFileName)

	if err != nil {
		log.Fatalf("%s", err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Lets make sure ec2 is available should be a Waiter refactor later
	time.Sleep(30 * time.Second)

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", addr, 22), sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil

}

// SendCommands a slice of commands for automated testing
func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	cmd := strings.Join(cmds, "; ")
	output, err := session.Output(cmd)
	if err != nil {
		return output, fmt.Errorf("failed to execute command '%s' on server: %v", cmd, err)
	}

	return output, err
}
