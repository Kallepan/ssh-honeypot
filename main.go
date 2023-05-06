package main

import (
	"fmt"

	"github.com/kallepan/ssh-honeypot/conf"
	"github.com/kallepan/ssh-honeypot/logger"
	"github.com/kallepan/ssh-honeypot/ssh"
)

func main() {
	conf.LoadEnv()
	logger.StartLogger()
	sshOpts := conf.GetOpts()

	logger.Info("Starting SSH honeypot")
	println(fmt.Sprintf("Starting SSH honeypot on port %d\n", sshOpts.Port))
	ssh.Listen(sshOpts)
}
