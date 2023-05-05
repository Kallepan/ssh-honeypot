package ssh

type SSHOpts struct {
	ServerCiphers []string
	ServerMACs    []string
	Port          int
}
