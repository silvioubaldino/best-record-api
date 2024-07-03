package sshclient

import (
	"bufio"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func NewSSHClient() (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("SSH_USERNAME"), // Replace with your SFTP username
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("SSH_PASSWORD")), // Replace with your SFTP password
		},
		HostKeyCallback: ssh.HostKeyCallback(callback),
	}

	sshClient, err := ssh.Dial("tcp", "localhost:22", sshConfig) // Replace with your SFTP server IP and port
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func callback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func readAuthorizedKeys() ([]ssh.PublicKey, error) {
	authorizedKeys := []ssh.PublicKey{}

	file, err := os.Open(".ssh/authorized_keys")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		publicKey, err := ssh.ParsePublicKey([]byte(line))
		if err != nil {
			return nil, err
		}

		authorizedKeys = append(authorizedKeys, publicKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return authorizedKeys, nil
}
