package main

import (
	"github.com/kallepan/ssh-honeypot/internal/config"
	"github.com/kallepan/ssh-honeypot/internal/ssh"
	"github.com/kallepan/ssh-honeypot/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Get SSH options
	sshOpts := config.SSHOpts{}
	sshOpts.Init()

	// Start SSH server
	ssh.Listen(sshOpts)
}
