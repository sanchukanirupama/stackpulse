# StackPulse Usage Guide

## Building the Binary

```bash
# Build the binary
make build

# Or build manually
go build -o build/stackpulse .
```

## Running StackPulse

### 1. Monitor a Node.js service by port (most common)

```bash
# Monitor a service running on port 3000
./build/stackpulse watch --port 3000

# With custom thresholds and faster polling
./build/stackpulse watch --port 3000 --cpu-threshold 80 --heap-limit 200MB --polling-ms 50 --inspect-port 9229
```

### 2. Monitor by Process ID (PID)

```bash
# If you know the PID directly
./build/stackpulse watch --pid 1234

# With custom settings
./build/stackpulse watch --pid 1234 --cpu-threshold 70 --polling-ms 100 --inspect-port 9229
```

## Quick Start Examples

### Example 1: Monitor Express.js App
```bash
# Start your Node.js app with inspector
node --inspect=0.0.0.0:9229 app.js

# In another terminal, monitor it
./build/stackpulse watch --port 3000 --polling-ms 100 --inspect-port 9229
```

### Example 2: High-Frequency Monitoring
```bash
# Monitor with millisecond precision for leak detection
./build/stackpulse watch --port 3000 --polling-ms 10 --cpu-threshold 50 --heap-limit 100MB --inspect-port 9229
```

### Example 3: Development Mode
```bash
# Quick development monitoring
make dev
# This runs: ./build/stackpulse watch --port 3000 --polling-ms 100 --inspect-port 9229
```

## Command Line Options

### Watch Command Options
- `--host`: Host to monitor (default: 127.0.0.1)
- `--port`: Port to monitor
- `--pid`: Process ID to monitor
- `--heap-limit`: Heap memory limit threshold (default: 150MB)
- `--cpu-threshold`: CPU usage threshold percentage (default: 70)
- `--polling-ms`: Polling interval in milliseconds (default: 100)
- `--inspect-port`: V8 inspector port (default: 9229)

## Troubleshooting

### Common Issues

1. **"No process listening on port"**
   ```bash
   # Check if your Node.js app is running
   lsof -i :3000
   
   # Or use netstat
   netstat -tlnp | grep :3000
   ```

2. **"Failed to connect to V8 inspector"**
   ```bash
   # Start your Node.js app with inspector enabled
   node --inspect=0.0.0.0:9229 app.js
  
  # Or specify a custom port and match it in stackpulse
  node --inspect=0.0.0.0:9230 app.js
  stackpulse watch --port 3000 --inspect-port 9230
   ```

3. **Permission denied**
   ```bash
   # Make binary executable
   chmod +x build/stackpulse
   ```

4. **High CPU usage from monitoring**
   ```bash
   # Increase polling interval
   ./build/stackpulse watch --port 3000 --polling-ms 1000
   ```

## Performance Tips

- Use `--polling-ms 100` for general monitoring
- Use `--polling-ms 10-50` for leak detection
- Use `--polling-ms 1000` for low-overhead monitoring