package conf

import "strconv"

type SSHOpts struct {
	ServerCiphers    []string
	ServerMACs       []string
	ServerAlgorithms []string
	Port             int
}

func GetOpts() SSHOpts {
	port, err := strconv.Atoi(GetValueFromEnv("SSH_PORT"))
	if err != nil {
		port = 2222
	}

	sshOpts := SSHOpts{
		Port: port,
		ServerMACs: []string{
			"hmac-sha2-256-etm@openssh.com",
			"hmac-sha2-256",
			"hmac-sha1",
			"hmac-sha1-96",
		},
		ServerCiphers: []string{
			"aes128-ctr",
			"aes192-ctr",
			"aes256-ctr",
		},
		ServerAlgorithms: []string{"rsa", "ecdsa", "ed25519"},
	}

	return sshOpts
}
