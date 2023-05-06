# SSH Honeypot

## Description

This is a simple SSH honeypot written in GoLang. It accepts any username/password and username/key combination. It logs authentication attempts to a file. Furthermore, a fake shell is provided to the attacker, which logs all commands to a file. Allowed commands can be configured in the config/cmds.txt file.

The project can be built inside a Docker container. The Dockerfile is provided.

## Usage

Configure the honeypot by editing the .env file.
