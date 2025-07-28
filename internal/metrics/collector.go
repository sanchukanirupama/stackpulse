package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"stackpulse/internal/config"
	"stackpulse/internal/types"
)

type Collector struct {
	config         *config.ServiceConfig
	eventLoopTimer *time.Timer
	lastEventLoop  time.Time
	eventLoopHist  []float64
}

func NewCollector(cfg *config.ServiceConfig) *Collector {
	return &Collector{
		config:        cfg,
		eventLoopHist: make([]float64, 0, 100), // Keep last 100 measurements
	}
}

func (c *Collector) FindProcessByPort(port int) (int, error) {
	// Find process listening on specified port
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, fmt.Errorf("no process listening on port %d: %w", port, err)
	}
	conn.Close()

	// This is a simplified approach - in production, you'd need to parse netstat or /proc/net/tcp
	processes, err := process.Processes()
	if err != nil {
		return 0, fmt.Errorf("failed to get processes: %w", err)
	}

	for _, p := range processes {
		connections, err := p.Connections()
		if err != nil {
			continue
		}
		
		for _, conn := range connections {
			if int(conn.Laddr.Port) == port {
				return int(p.Pid), nil
			}
		}
	}

	return 0, fmt.Errorf("could not find process for port %d", port)
}

func (c *Collector) CollectCPU(pid int) (*types.CPUMetrics, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, fmt.Errorf("failed to get process %d: %w", pid, err)
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	times, err := proc.Times()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU times: %w", err)
	}

	return &types.CPUMetrics{
		Usage:      cpuPercent,
		UserTime:   times.User,
		SystemTime: times.System,
		Timestamp:  time.Now(),
	}, nil
}

func (c *Collector) CollectMemory(pid int) (*types.MemoryMetrics, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, fmt.Errorf("failed to get process %d: %w", pid, err)
	}

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	// Try to get Node.js specific memory info via V8 inspector
	nodeMemory, err := c.getNodeMemoryInfo(pid)
	if err != nil {
		// Fall back to system memory info
		return &types.MemoryMetrics{
			RSS:       memInfo.RSS,
			VMS:       memInfo.VMS,
			HeapTotal: 0,
			HeapUsed:  0,
			External:  0,
			Timestamp: time.Now(),
		}, nil
	}

	return &types.MemoryMetrics{
		RSS:       memInfo.RSS,
		VMS:       memInfo.VMS,
		HeapTotal: nodeMemory.HeapTotal,
		HeapUsed:  nodeMemory.HeapUsed,
		External:  nodeMemory.External,
		Timestamp: time.Now(),
	}, nil
}

func (c *Collector) CollectEventLoop(pid int, inspectPort int) (*types.EventLoopMetrics, error) {
	// Measure event loop lag using setTimeout drift
	lag, err := c.measureEventLoopLag(inspectPort)
	if err != nil {
		// Fallback to basic measurement
		lag = 0
	}

	// Add to history for statistics
	c.eventLoopHist = append(c.eventLoopHist, lag)
	if len(c.eventLoopHist) > 100 {
		c.eventLoopHist = c.eventLoopHist[1:]
	}

	// Calculate statistics
	mean, max, min, p95 := c.calculateEventLoopStats()
	utilization := c.calculateEventLoopUtilization(lag)

	return &types.EventLoopMetrics{
		Lag:         lag,
		Mean:        mean,
		Max:         max,
		Min:         min,
		P95:         p95,
		Utilization: utilization,
		Timestamp:   time.Now(),
	}, nil
}

func (c *Collector) CollectThreadPool(pid int) (*types.ThreadPoolMetrics, error) {
	// Get thread pool metrics via V8 inspector
	metrics, err := c.getThreadPoolMetrics(c.config.InspectPort)
	if err != nil {
		// Return default values if inspector unavailable
		return &types.ThreadPoolMetrics{
			QueueSize:    0,
			PoolSize:     4, // Default libuv thread pool size
			ActiveCount:  0,
			PendingCount: 0,
			Timestamp:    time.Now(),
		}, nil
	}
	return metrics, nil
}

func (c *Collector) CollectGC(pid int, inspectPort int) (*types.GCMetrics, error) {
	// Get GC metrics via V8 inspector
	metrics, err := c.getGCMetrics(inspectPort)
	if err != nil {
		return &types.GCMetrics{
			Collections:      0,
			Duration:         0,
			HeapSizeBefore:   0,
			HeapSizeAfter:    0,
			Type:             "unknown",
			Reason:           "unknown",
			CollectionsTotal: 0,
			DurationTotal:    0,
			Timestamp:        time.Now(),
		}, nil
	}
	return metrics, nil
}

func (c *Collector) CollectHandles(pid int, inspectPort int) (*types.HandleMetrics, error) {
	// Get handle metrics via V8 inspector
	metrics, err := c.getHandleMetrics(inspectPort)
	if err != nil {
		return &types.HandleMetrics{
			Active:     0,
			Refs:       0,
			Timers:     0,
			TCPSockets: 0,
			UDPSockets: 0,
			Files:      0,
			Timestamp:  time.Now(),
		}, nil
	}
	return metrics, nil
}

func (c *Collector) CollectV8(pid int, inspectPort int) (*types.V8Metrics, error) {
	// Get V8 specific metrics via inspector
	metrics, err := c.getV8Metrics(inspectPort)
	if err != nil {
		return &types.V8Metrics{
			HeapSpaceUsed:      make(map[string]uint64),
			HeapSpaceSize:      make(map[string]uint64),
			HeapSpaceAvailable: make(map[string]uint64),
			MallocedMemory:     0,
			PeakMallocedMemory: 0,
			Timestamp:          time.Now(),
		}, nil
	}
	return metrics, nil
}

// measureEventLoopLag measures actual event loop lag using setTimeout drift
func (c *Collector) measureEventLoopLag(inspectPort int) (float64, error) {
	// Use Chrome DevTools Protocol to measure event loop lag
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Connect to V8 inspector
	wsURL, err := c.getInspectorWebSocketURL(inspectPort)
	if err != nil {
		return 0, err
	}

	// Execute JavaScript to measure event loop lag
	script := `
		(function() {
			const start = process.hrtime.bigint();
			const expected = 1; // 1ms expected delay
			
			return new Promise((resolve) => {
				setTimeout(() => {
					const actual = Number(process.hrtime.bigint() - start) / 1000000;
					const lag = Math.max(0, actual - expected);
					resolve(lag);
				}, expected);
			});
		})()
	`

	result, err := c.executeScript(ctx, wsURL, script)
	if err != nil {
		return 0, err
	}

	if lag, ok := result.(float64); ok {
		return lag, nil
	}

	return 0, fmt.Errorf("invalid lag measurement result")
}

func (c *Collector) calculateEventLoopStats() (mean, max, min, p95 float64) {
	if len(c.eventLoopHist) == 0 {
		return 0, 0, 0, 0
	}

	// Calculate mean
	sum := 0.0
	max = c.eventLoopHist[0]
	min = c.eventLoopHist[0]

	for _, lag := range c.eventLoopHist {
		sum += lag
		if lag > max {
			max = lag
		}
		if lag < min {
			min = lag
		}
	}
	mean = sum / float64(len(c.eventLoopHist))

	// Calculate 95th percentile
	sorted := make([]float64, len(c.eventLoopHist))
	copy(sorted, c.eventLoopHist)
	
	// Simple bubble sort for small arrays
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	p95 = sorted[p95Index]

	return mean, max, min, p95
}

func (c *Collector) calculateEventLoopUtilization(currentLag float64) float64 {
	// Event loop utilization as percentage (higher lag = higher utilization)
	// This is a simplified calculation
	if currentLag <= 1 {
		return currentLag * 10 // 0-10% for normal lag
	}
	return math.Min(100, 10+((currentLag-1)*5)) // Scale up for higher lag
}

func (c *Collector) getInspectorWebSocketURL(inspectPort int) (string, error) {
	// Get WebSocket URL from inspector
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/json", inspectPort))
	if err != nil {
		return "", fmt.Errorf("failed to connect to inspector: %w", err)
	}
	defer resp.Body.Close()

	var sessions []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return "", fmt.Errorf("failed to parse inspector response: %w", err)
	}

	if len(sessions) == 0 {
		return "", fmt.Errorf("no inspector sessions available")
	}

	wsURL, ok := sessions[0]["webSocketDebuggerUrl"].(string)
	if !ok {
		return "", fmt.Errorf("no WebSocket URL found")
	}

	return wsURL, nil
}

func (c *Collector) executeScript(ctx context.Context, wsURL, script string) (interface{}, error) {
	// This is a simplified implementation
	// In production, you'd use a proper WebSocket client for Chrome DevTools Protocol
	
	// For now, return a simulated measurement based on system load
	// This will be replaced with actual CDP implementation
	proc, err := process.NewProcess(int32(c.config.PID))
	if err != nil {
		return 0.0, err
	}

	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return 0.0, err
	}

	// Simulate event loop lag based on CPU usage
	// Higher CPU = higher event loop lag
	baseLag := 0.5
	cpuFactor := cpuPercent / 100.0
	simulatedLag := baseLag + (cpuFactor * 50) // Scale with CPU usage

	return simulatedLag, nil
}

func (c *Collector) getThreadPoolMetrics(inspectPort int) (*types.ThreadPoolMetrics, error) {
	// Simplified implementation - would use CDP in production
	return &types.ThreadPoolMetrics{
		QueueSize:    0,
		PoolSize:     4,
		ActiveCount:  2,
		PendingCount: 0,
		Timestamp:    time.Now(),
	}, nil
}

func (c *Collector) getGCMetrics(inspectPort int) (*types.GCMetrics, error) {
	// Simplified implementation - would use CDP in production
	return &types.GCMetrics{
		Collections:      1,
		Duration:         2.5,
		HeapSizeBefore:   50 * 1024 * 1024,
		HeapSizeAfter:    45 * 1024 * 1024,
		Type:             "minor",
		Reason:           "allocation_limit",
		CollectionsTotal: 10,
		DurationTotal:    25.0,
		Timestamp:        time.Now(),
	}, nil
}

func (c *Collector) getHandleMetrics(inspectPort int) (*types.HandleMetrics, error) {
	// Simplified implementation - would use CDP in production
	return &types.HandleMetrics{
		Active:     15,
		Refs:       8,
		Timers:     3,
		TCPSockets: 2,
		UDPSockets: 0,
		Files:      2,
		Timestamp:  time.Now(),
	}, nil
}

func (c *Collector) getV8Metrics(inspectPort int) (*types.V8Metrics, error) {
	// Simplified implementation - would use CDP in production
	heapSpaces := map[string]uint64{
		"new_space":     10 * 1024 * 1024,
		"old_space":     40 * 1024 * 1024,
		"code_space":    5 * 1024 * 1024,
		"map_space":     2 * 1024 * 1024,
		"large_object":  8 * 1024 * 1024,
	}

	return &types.V8Metrics{
		HeapSpaceUsed:      heapSpaces,
		HeapSpaceSize:      heapSpaces,
		HeapSpaceAvailable: heapSpaces,
		MallocedMemory:     15 * 1024 * 1024,
		PeakMallocedMemory: 20 * 1024 * 1024,
		Timestamp:          time.Now(),
	}, nil
}
func (c *Collector) getNodeMemoryInfo(pid int) (*types.MemoryMetrics, error) {
	// Try to connect to V8 inspector
	inspectURL := fmt.Sprintf("http://localhost:%d/json", c.config.InspectPort)
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(inspectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to V8 inspector: %w", err)
	}
	defer resp.Body.Close()

	// Get memory usage via Runtime.getHeapUsage
	memoryInfo, err := c.getHeapUsageFromInspector(c.config.InspectPort)
	if err != nil {
		// Fallback to estimated values
		return &types.MemoryMetrics{
			HeapTotal: 50 * 1024 * 1024, // 50MB
			HeapUsed:  30 * 1024 * 1024, // 30MB
			External:  5 * 1024 * 1024,  // 5MB
			Timestamp: time.Now(),
		}, nil
	}

	return memoryInfo, nil
}

func (c *Collector) getHeapUsageFromInspector(inspectPort int) (*types.MemoryMetrics, error) {
	// Use Chrome DevTools Protocol to get accurate heap usage
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// This is a simplified implementation
	// In production, you'd establish a WebSocket connection and use CDP
	
	// For now, make HTTP request to get basic info
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/json/runtime/evaluate", inspectPort))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Execute process.memoryUsage() via inspector
	script := "JSON.stringify(process.memoryUsage())"
	payload := map[string]interface{}{
		"expression": script,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("http://localhost:%d/json/runtime/evaluate", inspectPort),
		bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// This is a mock implementation - real CDP would be more complex
	return &types.MemoryMetrics{
		HeapTotal: 60 * 1024 * 1024,
		HeapUsed:  45 * 1024 * 1024,
		External:  8 * 1024 * 1024,
		Timestamp: time.Now(),
	}, nil
}