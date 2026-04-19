# portwatch

Lightweight CLI to monitor and alert on port state changes across hosts.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Monitor one or more hosts and get alerted when a port changes state:

```bash
# Watch a single host
portwatch --host example.com --port 443

# Watch multiple hosts with a custom interval
portwatch --host example.com --host internal.host --port 80,443 --interval 30s

# Output alerts to a log file
portwatch --host example.com --port 22 --log alerts.log
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | — | Host(s) to monitor (repeatable) |
| `--port` | — | Port(s) to watch (comma-separated) |
| `--interval` | `60s` | Poll interval |
| `--log` | stdout | File path for alert output |
| `--timeout` | `5s` | Connection timeout per check |

### Example Output

```
2024/01/15 10:32:01 [OPEN]   example.com:443
2024/01/15 10:33:01 [CLOSED] example.com:443  — state changed
2024/01/15 10:34:01 [OPEN]   example.com:443  — state changed
```

## License

MIT © 2024 portwatch contributors