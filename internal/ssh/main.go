package ssh

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kallepan/ssh-honeypot/internal/config"
	"github.com/kallepan/ssh-honeypot/internal/fakeshell"
	"github.com/kallepan/ssh-honeypot/pkg/logger"
	"golang.org/x/crypto/ssh"
)

// handle the incoming ssh connection and present a fake shell
func handleServerConn(user string, remoteAddr string, chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		// Only allow session channels
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", newChan.ChannelType()))
			continue
		}

		connection, requests, err := newChan.Accept()
		if err != nil {
			logger.Errorf("Could not accept channel: %v", err)
			continue
		}

		// Create fake shell
		fakeshell.FakeShell(connection, requests, user, remoteAddr)

		go func() {
			for req := range requests {
				switch req.Type {
				case "shell":
					// Only accept the default shell. Commands are ignored
					if len(req.Payload) == 0 {
						req.Reply(true, nil)
					}
				case "pty-req":
					req.Reply(true, nil)
				}
			}
		}()

	}
}

func listen(conf *ssh.ServerConfig, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("Could not listen on port %d: %v", port, err)
	}

	for {
		// Accept connections
		conn, err := listener.Accept()
		if err != nil {
			logger.Errorf("Could not accept connection: %v", err)
			continue
		}

		// Handshake
		sConn, chans, reqs, err := ssh.NewServerConn(conn, conf)
		if err != nil {
			logger.Errorf("Could not handshake: %v", err)
			continue
		}

		// Log connection
		logger.Infof("New connection from %s (%s) as %s authenticated with %s",
			sConn.RemoteAddr(),
			sConn.ClientVersion(),
			sConn.User(),
			sConn.Permissions.Extensions["auth-type"],
		)

		// Incoming requests
		go ssh.DiscardRequests(reqs) // These are not used
		go handleServerConn(sConn.User(), sConn.RemoteAddr().String(), chans)
	}
}

func Listen(opts config.SSHOpts) {
	// Listen for SSH connections
	conf := &ssh.ServerConfig{
		Config: ssh.Config{
			Ciphers: opts.ServerCiphers,
			MACs:    opts.ServerMACs,
		},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			// Simulate an accepted key
			logger.Infof("Accepted key from %s", conn.RemoteAddr())
			return &ssh.Permissions{Extensions: map[string]string{"auth-type": "key"}}, nil
		},
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Simulate an accepted password
			logger.Infof("Accepted password from %s", conn.RemoteAddr())
			return &ssh.Permissions{Extensions: map[string]string{"auth-type": "pass"}}, nil
		},
	}

	// Add host keys
	hostkeys, err := config.SetupHostKeys(opts.ServerAlgorithms, filepath.Join("assets", "keys"))
	if err != nil {
		logger.Fatalf("Could not setup host keys: %v", err)
	}

	for _, key := range hostkeys {
		conf.AddHostKey(key)
	}

	// Get port
	port := os.Getenv("SSH_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatalf("Could not convert port to integer: %v", err)
	}

	logger.Infof("Listening on port %d", portInt)
	// Finally listen for connections
	listen(conf, portInt)
}
