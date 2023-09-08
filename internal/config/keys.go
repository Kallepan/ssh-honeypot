package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kallepan/ssh-honeypot/pkg/logger"
	"golang.org/x/crypto/ssh"
)

func SetupHostKeys(algorithms []string, dir string) ([]ssh.Signer, error) {
	var hostkeys []ssh.Signer
	for _, algorithm := range algorithms {
		keypath := filepath.Join(dir, fmt.Sprintf("ssh_host_%s_key", algorithm))
		if _, err := os.Stat(keypath); err != nil {
			logger.Infof("Generating %s key", algorithm)
			// Keys do not exist, generate them
			_, err := exec.Command("ssh-keygen", "-t", algorithm, "-f", keypath, "-N", "").Output()
			if err != nil {
				logger.Errorf("Could not generate %s key: %v", algorithm, err)
				continue
			}

			logger.Errorf("Generated %s key in %s", algorithm, keypath)
		}

		keyData, err := os.ReadFile(keypath)
		if err != nil {
			logger.Errorf("Could not read %s key: %v", algorithm, err)
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			logger.Errorf("Could not parse %s key: %v", algorithm, err)
			return nil, err
		}

		hostkeys = append(hostkeys, signer)
	}

	return hostkeys, nil
}
