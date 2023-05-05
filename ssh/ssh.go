package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func cleanCommand(cmd string) string {
	// For now, allow all commands
	return cmd
}

func handleServerConn(authType string, chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChan.Accept()
		if err != nil {
			// TODO Handle Error
			continue
		}

		go func(in <-chan *ssh.Request) {
			defer channel.Close()
			for req := range in {
				payload := cleanCommand(string(req.Payload))

				// Log the request
				log.Panicf(authType, payload)

				// TODO Handle request
				channel.SendRequest("exit-status", true, []byte{0, 0, 0, 0})
			}
		}(requests)
	}
}

func listen(config *ssh.ServerConfig, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Accept connections
		conn, err := listener.Accept()
		if err != nil {
			// TODO Handle Error
			continue
		}

		// Handshake
		sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			// TODO Handle Error
			continue
		}

		// Incoming requests
		go ssh.DiscardRequests(reqs)
		go handleServerConn(sConn.Permissions.Extensions["auth-type"], chans)
	}
}

func Listen(opts SSHOpts) {
	// Listen for SSH connections
	config := &ssh.ServerConfig{
		Config: ssh.Config{
			Ciphers: opts.ServerCiphers,
			MACs:    opts.ServerMACs,
		},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			// Simulate an accepted key
			return &ssh.Permissions{Extensions: map[string]string{"auth-type": "key-0"}}, nil
		},
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Simulate an accepted password
			return &ssh.Permissions{Extensions: map[string]string{"auth-type": "pass-0"}}, nil
		},
		NoClientAuth: true,
		NoClientAuthCallback: func(conn ssh.ConnMetadata) (*ssh.Permissions, error) {
			// Simulate an accepted key
			return &ssh.Permissions{Extensions: map[string]string{"auth-type": "anon-0"}}, nil
		},
	}

	// Add host key
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory")
	}
	keyPath := filepath.Join(wd, "ssh/id_rsa")
	if _, err := os.Stat(keyPath); err != nil {
		os.Mkdir(filepath.Join(wd, "ssh"), os.ModePerm)

		_, err := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", keyPath, "-N", "").Output()
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to generate key: %s", err))
		}
		fmt.Println(fmt.Sprintf("Generated key in %s", keyPath))
	}

	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("Failed to load private key")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Finally listen for connections
	go listen(config, opts.Port)
}
