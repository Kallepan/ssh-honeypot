# SSH Honeypot

## Description

This is a simple SSH honeypot written in GoLang. It accepts any username/password and username/key combination. It logs authentication attempts to a file. Furthermore, a fake shell is provided to the attacker, which logs all commands to a file. Allowed commands can be configured in the config/cmds.txt file.

The project can be built inside a Docker container. The Dockerfile is provided.

## Usage

### Development

```bash
# Configure the honeypot by editing the .env.example file.
cp .env.example .dev.env

# Run the honeypot
bash scripts/launch.dev.sh
```

### Production

```bash
# Configure the honeypot by editing the .env.example file.
cp .env.example .env

# Run the honeypot
bash scripts/launch.sh
```
