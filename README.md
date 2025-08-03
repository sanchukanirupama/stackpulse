# StackPulse - Node.js Memory & CPU Monitoring Tool

A lightweight monitoring tool for Node.js applications that tracks CPU usage, memory consumption, and performance metrics in real-time.

## Features

- **Real-time Monitoring**: Tracks CPU, memory, event loop lag, and thread pool metrics
- **Process Discovery**: Find Node.js processes by port or PID
- **Memory Leak Detection**: Advanced heap monitoring and garbage collection analysis
- **Event Loop Monitoring**: Detects and reports event loop blockages and lag
- **Live Dashboard**: Real-time terminal dashboard with color-coded alerts
- **Threshold Alerts**: Configurable alerting for CPU and memory thresholds

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/sanchukanirupama/stackpulse.git
cd stackpulse

# Build the binary
make build

# Install globally (optional)
make install
```

### Homebrew (macOS/Linux)

```bash
# Coming soon
```

### Linux Package Managers

```bash
# Coming soon
```

## Quick Start

### Monitor a Node.js service by port
```bash
stackpulse watch --port 3000 --polling-ms 100
```

### Monitor by process ID
```bash
stackpulse watch --pid 1234 --cpu-threshold 80 --heap-limit 200MB
```

## Usage

### Starting Your Node.js Application

For best monitoring results, start your Node.js application with the V8 inspector enabled:

```bash
# Enable V8 inspector on all interfaces
node --inspect=0.0.0.0:9229 app.js
```

### Watch Command

Monitor a Node.js process in real-time:

```bash
stackpulse watch [flags]

Flags:
  --host string          Host to monitor (default "127.0.0.1")
  --port int             Port to monitor
  --pid int              Process ID to monitor
  --heap-limit string    Heap memory limit threshold (default "150MB")
  --cpu-threshold float  CPU usage threshold percentage (default 70)
  --polling-ms int       Polling interval in milliseconds (default 100)
  --inspect-port int     V8 inspector port (default 9229)
```

## Examples

### Monitor Express.js App
```bash
# Start your Node.js app with inspector
node --inspect=0.0.0.0:9229 app.js

# Or use a custom port
node --inspect=0.0.0.0:9230 app.js

# Then monitor with matching inspect port
stackpulse watch --port 3000 --inspect-port 9230
# In another terminal, monitor it
stackpulse watch --port 3000 --polling-ms 100 --inspect-port 9229
```

### High-Frequency Monitoring
```bash
# Monitor with millisecond precision for leak detection
stackpulse watch --port 3000 --polling-ms 10 --cpu-threshold 50 --heap-limit 100MB --inspect-port 9229
```

### Monitor Specific Process
```bash
# Monitor by PID with custom thresholds
stackpulse watch --pid 1234 --cpu-threshold 80 --heap-limit 200MB --inspect-port 9229
```

## Metrics Tracked

### Memory Metrics
- **RSS**: Resident Set Size (physical memory usage)
- **VMS**: Virtual Memory Size
- **Heap Total**: Total heap size allocated by V8
- **Heap Used**: Amount of heap currently in use
- **External**: Memory used by C++ objects bound to JavaScript

### CPU Metrics
- **Usage Percentage**: Current CPU utilization
- **User Time**: Time spent in user mode
- **System Time**: Time spent in kernel mode

### Event Loop Metrics
- **Lag**: Current event loop delay
- **Mean**: Average event loop delay
- **95th Percentile**: 95% of measurements below this value
- **Utilization**: Event loop utilization percentage

### Node.js Specific Metrics
- **Garbage Collection**: GC frequency and duration
- **Handle Count**: Active handles (timers, sockets, files)
- **Thread Pool**: Queue size and active threads
- **V8 Heap Spaces**: Detailed heap space usage

## Alerting

StackPulse provides threshold-based alerting for:

- CPU usage exceeding configured limits
- Memory usage approaching heap limits
- Event loop lag indicating performance issues
- Heap usage percentage thresholds

Alerts are displayed in the terminal dashboard with color-coded severity levels.

## Performance Tips

- Use `--polling-ms 100` for general monitoring
- Use `--polling-ms 10-50` for leak detection
- Use `--polling-ms 1000` for low-overhead monitoring

## System Requirements

- Go 1.21 or higher (for building from source)
- Linux, macOS, or Windows
- Node.js applications with V8 inspector access (recommended)

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for development guidelines.

## License

MIT License - see LICENSE file for details.
