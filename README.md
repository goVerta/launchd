# launchd
Launchd is a simple tool that helps you deploy a Go program (a single binary) to a Linux server that uses systemd. Think of it as a very small CI/CD helper you can run from your laptop.

It does five things for you:
- Build your Go app locally.
- Copy the binary to your server over SSH.
- Create/refresh a systemd service for it.
- Optionally run database migrations.
- Check that the app is healthy before finishing.

If you’re new to servers, follow this guide step by step. You don’t need Docker or a complex CI system.

## What you need (prerequisites)
- A Go project with a `main` package you can build.
- Your laptop/workstation with:
  - Go installed
  - OpenSSH client (`ssh`, `scp`)
- A Linux server with:
  - systemd (most modern distros have it)
  - SSH access (you can log in with a user that can use `sudo`)
  - Network port open for your app (e.g., 8080)
- Your app should expose a health endpoint like `GET /health` that returns 200 OK when ready.

## Install / Build launchd
Clone this repo and build the CLI:

```bash
git clone <this-repo-url>
cd launchd
go build -o launchd ./cmd/launchd
```

This creates a `launchd` binary in the project root. Optionally, move it into your PATH:

```bash
mv launchd /usr/local/bin/
```

## Prepare your app
Make sure your app builds locally:

```bash
go build -o ./bin/myapp ./path/to/your/cmd
```

Also ensure your app can start with a port flag or env (for example `--port 8080`) and serves `/health`.

## Quick start: deploy in one command
Run from your laptop (replace placeholders):

```bash
launchd deploy \
  --host <server-ip-or-dns> \
  --user <ssh-username> \
  --app  ./path/to/your/app \
  --port <port> \
  --timeout 60s
```

What happens:
1) Your app is compiled.
2) The binary is copied to `/usr/local/bin/<app>` on the server.
3) A systemd unit `<app>.service` is created/updated and (re)started.
4) Optional migrations run if configured.
5) Launchd waits for `http://<host>:<port>/health` to be OK.

On success, your service is enabled on boot and running under systemd.

## How it works (simple explanation)
- Build: uses your local Go toolchain to compile your app.
- Transfer: uses `scp` over SSH to put the binary on the server.
- Service: writes a `.service` file, runs `systemctl daemon-reload`, `enable`, and `restart`.
- Migrate (optional): calls a migration tool if you have one installed on the server.
- Health check: polls your `/health` endpoint until it responds OK or times out.

## Common problems and fixes
- SSH fails: verify `ssh <user>@<host>` works and keys/Passwords are set up.
- Sudo prompts: ensure your SSH user can run the necessary `systemctl` and file copy with `sudo`.
- Port in use: stop whatever is using that port or change `--port`.
- Health check fails: make sure your service starts quickly, listens on the right port, and returns 200 on `/health`.
- See logs: `ssh <user>@<host> 'sudo journalctl -u <app>.service -f'`.

## Remove or stop the service (on the server)
```bash
sudo systemctl stop <app>.service
sudo systemctl disable <app>.service
sudo rm -f /usr/local/bin/<app>
sudo rm -f /etc/systemd/system/<app>.service
sudo systemctl daemon-reload
```

## Safety and re-runs
You can run the same deploy command again. It will overwrite the binary, refresh the unit, and restart safely. This is called idempotent behavior.

## FAQ
- Do I need Docker? No.
- Do I need Go on the server? No, only on your laptop (build happens locally).
- Do I need root? You need `sudo` to install the binary and manage systemd.

Built with care by Saad and the team.
