package fakeshell

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/kallepan/ssh-honeypot/conf"
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
	pathToCmds := conf.GetValueFromEnv("PATH_TO_CMDS")
	host := conf.GetValueFromEnv("SSH_HOST")
	if pathToCmds == "" {
		pathToCmds = "conf/cmds.txt"
	}

	bytes, err := ioutil.ReadFile(pathToCmds)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not read file: %s", pathToCmds))
	}
	commandsList := strings.Split(string(bytes), "\n")

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

		for _, cmd := range commandsList {
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
