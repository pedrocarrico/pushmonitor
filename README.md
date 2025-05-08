# Pushmonitor

**Pushmonitor** is a lightweight, Go-based service that performs push monitoring tests. It is designed to run as a `systemd` service and is compatible with tools like [Paessler PRTG](https://www.paessler.com/support/how-to/http-push-monitoring), [PushMon](https://www.pushmon.com/), and [StatusCake](https://www.statuscake.com/kb/knowledge-base/what-is-push-monitoring/).

## Features

- Executes push monitoring tests at configurable intervals
- Supports multiple independent test configurations
- Optionally runs a system command before performing a test
- Handles retries and HTTP timeouts
- Simple YAML-based configuration
- Logs to a file with configurable verbosity
- Designed to run reliably under `systemd`

## Installation

1. **Build the binary:**

```bash
make build
```

2. **Run the application:**

```bash
./pushmonitor
```

## Configuration

Pushmonitor looks for a configuration file in the following locations (in order):

1. `/etc/pushmonitor/config.yaml`
2. `config/config.yaml` (relative to the binary)

### Configuration options

- Multiple named push tests
- Custom test intervals and retry limits
- Optional pre-test system commands
- Logging options
- Global timeout setting
- Optional PID file path

### Example `config.yaml`

```yaml
push_tests:
  - name: "Heartbeat"
    url: "https://push.statuscake.com/?PK=12345&TestID=67890&time=0"
    interval: 300      # in seconds
    retries: 3

  - name: "Nginx status"
    url: "http://pshmn.com/ebFnY1"
    interval: 3600     # in seconds
    retries: 3
    command: "service nginx status"

logging:
  file: "/var/log/pushmonitor.log"
  level: "info"

pid_file: "/etc/pushmonitor/pid"
timeout: 30             # in seconds
```

### Running as a systemd Service

Create a service file at /etc/systemd/system/pushmonitor.service:

```
[Unit]
Description=Pushmonitor Service
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/etc/pushmonitor
ExecStart=/usr/local/bin/pushmonitor
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

- Start the service:
```bash
sudo systemctl start push-monitor
```

- Check service status:
```bash
sudo systemctl status push-monitor
```

- View logs:
```bash
sudo journalctl -u push-monitor
```

- Stop the service:
```bash
sudo systemctl stop push-monitor
```

- Restart the service:
```bash
sudo systemctl restart push-monitor
```

## Supported Platforms

Pushmonitor has been tested on:
- Amazon Linux 2023 (AMI)
- Debian 12.10 "Bookworm"

It is expected to work on other modern Linux distributions that support systemd.

## License

MIT License. See [LICENSE.md](./LICENSE.md) for details.

##Contributions

Contributions, issues, and feature requests are welcome! Please open a pull request or issue on GitHub.
