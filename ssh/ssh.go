package ssh

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kallepan/ssh-honeypot/conf"
	"golang.org/x/crypto/ssh"
)

var ErrorLog = conf.ErrorLog
var Info = conf.Log

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
			ErrorLog.Printf("Could not accept channel: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			defer func() {
				_ = channel.Close()
			}()
			for req := range in {
				payload := cleanCommand(string(req.Payload))

				// Log the request
				Info.Println(authType, payload)

				// TODO Handle request

				channel.SendRequest("exit-status", true, []byte{0, 0, 0, 0})
			}
		}(requests)
	}
}

func listen(config *ssh.ServerConfig, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		ErrorLog.Fatal(err)
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

func setupHostKeys(algorithms []string, dir string) ([]ssh.Signer, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var hostkeys []ssh.Signer
	for _, algorithm := range algorithms {
		keypath := filepath.Join(dir, fmt.Sprintf("ssh_host_%s_key", algorithm))
		if _, err := os.Stat(keypath); err != nil {
			Info.Printf("Generating %s key", algorithm)
			// Keys do not exist, generate them
			_, err := exec.Command("ssh-keygen", "-t", algorithm, "-f", keypath, "-N", "").Output()
			if err != nil {
				ErrorLog.Printf("Could not generate %s key: %v", algorithm, err)
				continue
			}

			ErrorLog.Printf("Generated %s key in %s", algorithm, keypath)
		}

		keyData, err := os.ReadFile(keypath)
		if err != nil {
			ErrorLog.Fatalf("Could not read %s key: %v", algorithm, err)
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			ErrorLog.Fatalf("Could not parse %s key: %v", algorithm, err)
			return nil, err
		}

		hostkeys = append(hostkeys, signer)
	}

	return hostkeys, nil
}

func Listen(opts conf.SSHOpts) {
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
		ErrorLog.Fatal("Could not get working directory")
	}
	hostkeys, err := setupHostKeys(opts.ServerAlgorithms, filepath.Join(wd, "keys"))

	if err != nil {
		ErrorLog.Fatal(fmt.Sprintf("Could not setup host keys: %v", err))
	}

	for _, key := range hostkeys {
		config.AddHostKey(key)
	}

	// Finally listen for connections
	go listen(config, opts.Port)
}
