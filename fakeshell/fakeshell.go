package fakeshell

import (
	"fmt"
	"io"
	"strings"

	"github.com/kallepan/ssh-honeypot/logger"
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
	term := term.NewTerminal(s, fmt.Sprintf(
		"%s%s@%s>%s ",
		"\x1b[0m", // green
		user,
		"localhost",
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

		if unknown {
			Write(term, fmt.Sprintf("bash: %s: command not found\n", command))
			Write(term, "1")
		}
	}

	_, err := term.Write([]byte("logout\n"))
	if err != nil {
		// Does not matter still close the connection
		s.Close()
		return
	}

	s.Close()
}
