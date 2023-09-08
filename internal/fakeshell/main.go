package fakeshell

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kallepan/ssh-honeypot/pkg/logger"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Write(w io.Writer, str string) {
	chars := strings.Split(str, "")

	for _, char := range chars {
		fmt.Fprint(w, char)
	}
}

// create a fake shell where the ssh user can "execute" commands
func FakeShell(s ssh.Channel, reqs <-chan *ssh.Request, user string, remoteAddr string) {
	cmdFile := os.Getenv("CMDS_FILE")
	host := os.Getenv("SSH_HOST")

	// read commands from file
	bytes, err := os.ReadFile(cmdFile)
	if err != nil {
		logger.Fatalf("Could not read file: %s", cmdFile)
	}
	commands := strings.Split(string(bytes), "\n")

	// create terminal
	term := term.NewTerminal(s, fmt.Sprintf(
		"%s%s@%s:~#%s ",
		"\x1b[0m", // green
		user,
		host,
		"\x1b[0m",
	))

	for {
		ln, err := term.ReadLine()
		if err != nil {
			break
		}
		if ln == "exit" {
			break
		}

		logger.Infof("Host: %s, User %s executed command: %s", remoteAddr, user, ln)

		commandAndArgs := strings.Split(ln, " ")
		command := commandAndArgs[0]
		unknown := true

		for _, cmd := range commands {
			if cmd == command {
				unknown = false
				break
			}
		}

		if unknown {
			Write(term, fmt.Sprintf("bash: %s: command not found\n", command))
		}
	}

	_, err = term.Write([]byte("logout\n"))
	if err != nil {
		s.Close()
		return
	}
	s.Close()
}
