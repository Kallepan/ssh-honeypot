package main

import (
	"github.com/kallepan/ssh-honeypot/conf"
	"github.com/kallepan/ssh-honeypot/ssh"
)

func main() {
	conf.StartLogger()
	sshOpts := conf.GetOpts()

	ssh.Listen(sshOpts)
}
