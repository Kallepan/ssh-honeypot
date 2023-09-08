package config

type SSHOpts struct {
	ServerCiphers    []string
	ServerMACs       []string
	ServerAlgorithms []string
}

func (s *SSHOpts) Init() {
	s.ServerMACs = []string{
		"hmac-sha2-256-etm@openssh.com",
		"hmac-sha2-256",
		"hmac-sha1",
		"hmac-sha1-96",
	}
	s.ServerCiphers = []string{
		"aes128-ctr",
		"aes192-ctr",
		"aes256-ctr",
	}
	s.ServerAlgorithms = []string{"rsa", "ecdsa", "ed25519"}
}
