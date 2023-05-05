package conf

type SSHOpts struct {
	ServerCiphers    []string
	ServerMACs       []string
	ServerAlgorithms []string
	Port             int
}

func GetOpts() SSHOpts {
	sshOpts := SSHOpts{
		Port: 2222,
		ServerMACs: []string{
			"",
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
