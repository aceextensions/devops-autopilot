# 🚀 DevOps Autopilot

A lightweight, CLI-first DevOps monitoring and alerting tool built in Go. Monitor server health, Docker containers, send alerts, and generate reports — all from a single binary. Runs as a native system service on **Linux**, **macOS**, and **Windows**.

---

## ✨ Features

- **System Health** — CPU, Memory, Disk, Load Average, Uptime
- **Docker Monitoring** — Running/stopped containers, restart counts, CPU & memory per container
- **Alert Engine** — Configurable thresholds with Slack & Email notifications
- **Report Generator** — Structured ASCII reports on demand
- **Continuous Monitor** — Scheduled polling loop with graceful shutdown
- **System Service** — Install as `systemd` (Linux), `launchd` (macOS), or Windows SCM service
- **Single Binary** — Compiles to one self-contained executable (<50MB RAM)

---

## 📋 Requirements

| Dependency | Notes |
|---|---|
| Go 1.21+ | For building from source |
| Docker (optional) | For Docker monitoring features |

---

## 🛠️ Installation

### Build from Source

```bash
git clone https://github.com/your-org/devops-autopilot.git
cd devops-autopilot
go mod tidy
go build -o devops-autopilot .
```

### Cross-compile

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o devops-autopilot-linux .

# Linux arm64 (Raspberry Pi, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o devops-autopilot-arm64 .

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o devops-autopilot-mac .

# Windows
GOOS=windows GOARCH=amd64 go build -o devops-autopilot.exe .
```

---

## ⚙️ Configuration

Edit `configs/config.yaml` before running:

```yaml
server_name: prod-01          # Label shown in reports and alerts

alerts:
  cpu: 80                     # Alert when CPU > 80%
  memory: 85                  # Alert when Memory > 85%
  disk: 90                    # Alert when Disk > 90%

slack:
  enabled: true
  webhook_url: https://hooks.slack.com/services/YOUR/WEBHOOK/URL

email:
  enabled: true
  smtp_host: smtp.gmail.com
  smtp_port: 587
  username: you@gmail.com
  password: your_app_password  # Use Gmail App Password
  to: admin@yourcompany.com

monitor:
  interval_seconds: 60        # How often to poll in monitor/service mode
```

### Environment Variable Overrides

Sensitive values can be overridden via environment variables (take priority over config file):

| Variable | Overrides |
|---|---|
| `SERVER_NAME` | `server_name` |
| `SLACK_WEBHOOK_URL` | `slack.webhook_url` |
| `SMTP_PASSWORD` | `email.password` |

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/..."
export SMTP_PASSWORD="my-app-password"
./devops-autopilot monitor
```

---

## 🖥️ CLI Commands

```
devops-autopilot [command]
```

### `health` — One-shot system check

```bash
./devops-autopilot health
```

```
Server   : prod-01
Uptime   : 5d 3h 42m
CPU      : 23.4%
Memory   : 61.3% (used 9.8GB / total 16.0GB)
Disk     : 62.7% (used 146.5GB / total 233.6GB)
Load Avg : 2.10 / 1.87 / 1.65
```

---

### `docker` — Docker container status

```bash
./devops-autopilot docker
```

```
Docker Status  : Available
Running        : 5
Stopped        : 2
Total          : 7

●  930b53624ed0  my-app    nginx:alpine   running    0
●  47c78b171f10  my-db     mysql:8.0      running    0
○  1139aa42886b  old-api   node:18        exited     3
```

> **●** = running, **○** = stopped/exited. `restarts` count helps identify crash-looping containers.

---

### `report` — Full report (system + docker + alerts)

```bash
./devops-autopilot report
```

```
╔══════════════════════════════════════╗
║  DevOps Autopilot Report             ║
╚══════════════════════════════════════╝

Server   : prod-01
Time     : 2026-03-03 20:45:00
Uptime   : 5d 3h 42m

── System Health ─────────────────────
CPU      : 23.4%
Memory   : 61.3% (9.8GB / 16.0GB)
Disk     : 62.7% (146.5GB / 233.6GB)
Load Avg : 2.10 / 1.87 / 1.65

── Docker ─────────────────────────────
Running  : 5
Stopped  : 2
Total    : 7

── Alerts ─────────────────────────────
None
```

With a custom config path:
```bash
./devops-autopilot --config /etc/devops-autopilot/config.yaml report
```

---

### `monitor` — Continuous polling loop

```bash
./devops-autopilot monitor
```

Runs a check every `interval_seconds` (default: 60s). Sends Slack/Email alerts when thresholds are exceeded. Press **Ctrl+C** for graceful shutdown.

---

### `service` — Native OS service management

> **Requires elevated privileges:** `sudo` on Linux/macOS, Administrator on Windows.

```bash
# Install and start the service
sudo ./devops-autopilot service install
sudo ./devops-autopilot service start

# Check status
sudo ./devops-autopilot service status

# Stop and remove
sudo ./devops-autopilot service stop
sudo ./devops-autopilot service uninstall
```

| Action | Description |
|---|---|
| `install` | Register with the OS init system |
| `uninstall` | Remove the registration |
| `start` | Start the background daemon |
| `stop` | Stop the daemon |
| `restart` | Restart the daemon |
| `status` | Show current service state |

---

## 🔧 How It Works

```
┌─────────────────────────────────────────────────┐
│                  CLI Layer (Cobra)               │
│   health │ docker │ report │ monitor │ service   │
└───────────────────────┬─────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────┐
│               Service Manager                    │
│      (kardianos/service → systemd/launchd/SCM)  │
└───────────────────────┬─────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
┌──────────────┐ ┌─────────────┐ ┌────────────────┐
│   Collector  │ │Alert Engine │ │  Notifier      │
│  system.go   │ │  engine.go  │ │  slack.go      │
│  docker.go   │ │             │ │  email.go      │
└──────────────┘ └─────────────┘ └────────────────┘
        │
        ▼
┌──────────────┐
│    Report    │
│  builder.go  │
└──────────────┘
```

### Component Breakdown

| Package | Role |
|---|---|
| `internal/config` | Loads `config.yaml`, applies env var overrides |
| `internal/collector/system.go` | Reads CPU/Memory/Disk/Load via `gopsutil` |
| `internal/collector/docker.go` | Queries Docker daemon socket for container states |
| `internal/alert/engine.go` | Compares metrics to thresholds, produces alert list |
| `internal/notifier/slack.go` | HTTP POST to Slack webhook |
| `internal/notifier/email.go` | SMTP email dispatch |
| `internal/report/builder.go` | Formats all data into readable ASCII report |
| `internal/service/manager.go` | Wraps daemon lifecycle using `kardianos/service` |

---

## 🐧 Linux — systemd Setup

After `service install`, a unit file is automatically created. You can inspect it:

```bash
sudo systemctl status devops-autopilot
sudo journalctl -u devops-autopilot -f    # Live logs
```

To use a custom config path, edit the unit file:

```bash
sudo systemctl edit devops-autopilot
```

Add:
```ini
[Service]
ExecStart=
ExecStart=/usr/local/bin/devops-autopilot --config /etc/devops-autopilot/config.yaml service run
```

---

## 🍎 macOS — launchd Setup

After `service install`, a plist is registered under `LaunchDaemons`. Logs go to:

```bash
tail -f /var/log/devops-autopilot.log
```

---

## 🪟 Windows — Service Control Manager

Run PowerShell as Administrator:

```powershell
.\devops-autopilot.exe service install
.\devops-autopilot.exe service start
```

View in Services panel (`services.msc`) as **DevOps Autopilot**.

---

## 📁 Project Structure

```
devops-autopilot/
├── main.go                      # Entry point
├── configs/
│   └── config.yaml              # Configuration file
├── cmd/
│   ├── root.go                  # Cobra root command
│   ├── health.go                # health subcommand
│   ├── docker.go                # docker subcommand
│   ├── report.go                # report subcommand
│   ├── monitor.go               # monitor subcommand
│   └── service.go               # service subcommand
└── internal/
    ├── config/loader.go         # Config loading
    ├── collector/
    │   ├── system.go            # System metrics
    │   └── docker.go            # Docker metrics
    ├── alert/engine.go          # Threshold evaluation
    ├── notifier/
    │   ├── slack.go             # Slack alerts
    │   └── email.go             # Email alerts
    ├── report/builder.go        # Report formatting
    └── service/manager.go       # OS service lifecycle
```

---

## 🔒 Security Best Practices

- **Never commit** `config.yaml` with real credentials. Use environment variables.
- Set restrictive file permissions: `chmod 600 configs/config.yaml`
- Use [Gmail App Passwords](https://support.google.com/accounts/answer/185833) instead of your main password.
- The binary requires Docker socket access (`/var/run/docker.sock`) — add your service user to the `docker` group on Linux.

---

## 🗺️ Roadmap

| Version | Feature |
|---|---|
| **v0.1** | ✅ CLI tool, systemd/launchd/Windows service, Slack/Email alerts |
| **v0.2** | Web dashboard (Fiber), REST API, JSON export, Prometheus endpoint |
| **v0.3** | Multi-server agent mode, central aggregator |
| **v1.0** | SaaS platform |

---

## 📄 License

MIT License — use freely, contributions welcome.
