package main

import (
	"fmt"
	"os"

	"github.com/kallepan/ssh-honeypot/conf"
	"github.com/kallepan/ssh-honeypot/logger"
	"github.com/kallepan/ssh-honeypot/ssh"
)

func main() {
	logger.StartLogger()
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--prod" {
			conf.LoadEnv()
		}
	}

	sshOpts := conf.GetOpts()

	logger.Info("Starting SSH honeypot")
	println(fmt.Sprintf("Starting SSH honeypot on port %d\n", sshOpts.Port))
	ssh.Listen(sshOpts)

}
