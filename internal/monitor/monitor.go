package monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"stackpulse/internal/config"
	"stackpulse/internal/metrics"
	"stackpulse/internal/display"
	"stackpulse/internal/alerts"
	"stackpulse/internal/types"
)

type Monitor struct {
	config     *config.ServiceConfig
	metrics    *metrics.Collector
	display    *display.Dashboard
	alerts     *alerts.Manager
	running    bool
	mu         sync.RWMutex
}

func New(cfg *config.ServiceConfig) *Monitor {
	return &Monitor{
		config:  cfg,
		metrics: metrics.NewCollector(cfg),
		display: display.NewDashboard(),
		alerts:  alerts.NewManager(),
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("monitor is already running")
	}
	m.running = true
	m.mu.Unlock()

	log.Printf("Starting monitor for PID: %d, Host: %s, Port: %d", 
		m.config.PID, m.config.Host, m.config.Port)

	ticker := time.NewTicker(m.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.mu.Lock()
			m.running = false
			m.mu.Unlock()
			log.Println("Monitor stopped")
			return nil
		case <-ticker.C:
			if err := m.collectAndProcess(); err != nil {
				log.Printf("Failed to collect metrics: %v", err)
			}
		}
	}
}

func (m *Monitor) collectAndProcess() error {
	// Get PID if not specified
	if m.config.PID == 0 {
		pid, err := m.metrics.FindProcessByPort(m.config.Port)
		if err != nil {
			return fmt.Errorf("failed to find process: %w", err)
		}
		m.config.PID = pid
	}

	// Collect all metrics
	cpuMetrics, err := m.metrics.CollectCPU(m.config.PID)
	if err != nil {
		return fmt.Errorf("failed to collect CPU metrics: %w", err)
	}

	memoryMetrics, err := m.metrics.CollectMemory(m.config.PID)
	if err != nil {
		return fmt.Errorf("failed to collect memory metrics: %w", err)
	}

	eventLoopMetrics, err := m.metrics.CollectEventLoop(m.config.PID, m.config.InspectPort)
	if err != nil {
		log.Printf("Warning: Failed to collect event loop metrics: %v", err)
		// Use default values
		eventLoopMetrics = &types.EventLoopMetrics{
			Lag:         0,
			Mean:        0,
			Max:         0,
			Min:         0,
			P95:         0,
			Utilization: 0,
			Timestamp:   time.Now(),
		}
	}

	threadPoolMetrics, err := m.metrics.CollectThreadPool(m.config.PID)
	if err != nil {
		log.Printf("Warning: Failed to collect thread pool metrics: %v", err)
		threadPoolMetrics = &types.ThreadPoolMetrics{
			QueueSize:    0,
			PoolSize:     4,
			ActiveCount:  0,
			PendingCount: 0,
			Timestamp:    time.Now(),
		}
	}

	// Collect additional Node.js specific metrics
	gcMetrics, err := m.metrics.CollectGC(m.config.PID, m.config.InspectPort)
	if err != nil {
		log.Printf("Warning: Failed to collect GC metrics: %v", err)
		gcMetrics = &types.GCMetrics{
			Collections:      0,
			Duration:         0,
			HeapSizeBefore:   0,
			HeapSizeAfter:    0,
			Type:             "unknown",
			Reason:           "unknown",
			CollectionsTotal: 0,
			DurationTotal:    0,
			Timestamp:        time.Now(),
		}
	}

	handleMetrics, err := m.metrics.CollectHandles(m.config.PID, m.config.InspectPort)
	if err != nil {
		log.Printf("Warning: Failed to collect handle metrics: %v", err)
		handleMetrics = &types.HandleMetrics{
			Active:     0,
			Refs:       0,
			Timers:     0,
			TCPSockets: 0,
			UDPSockets: 0,
			Files:      0,
			Timestamp:  time.Now(),
		}
	}

	v8Metrics, err := m.metrics.CollectV8(m.config.PID, m.config.InspectPort)
	if err != nil {
		log.Printf("Warning: Failed to collect V8 metrics: %v", err)
		v8Metrics = &types.V8Metrics{
			HeapSpaceUsed:      make(map[string]uint64),
			HeapSpaceSize:      make(map[string]uint64),
			HeapSpaceAvailable: make(map[string]uint64),
			MallocedMemory:     0,
			PeakMallocedMemory: 0,
			Timestamp:          time.Now(),
		}
	}

	// Create status
	status := &types.Status{
		PID:         m.config.PID,
		CPU:         *cpuMetrics,
		Memory:      *memoryMetrics,
		EventLoop:   *eventLoopMetrics,
		ThreadPool:  *threadPoolMetrics,
		GC:          *gcMetrics,
		Handles:     *handleMetrics,
		V8:          *v8Metrics,
		Timestamp:   time.Now(),
	}

	// Check for alerts
	alertList := m.alerts.CheckThresholds(status, m.config)
	status.Alerts = alertList

	// Update display
	m.display.Update(status)

	// Send alerts if any
	if len(alertList) > 0 {
		for _, alert := range alertList {
			log.Printf("ALERT [%s] %s: %s (Value: %.2f, Threshold: %.2f)", 
				string(alert.Severity), string(alert.Type), alert.Message, 
				alert.Value, alert.Threshold)
		}
	}

	return nil
}

func GetCurrentStatus() (*types.Status, error) {
	// Implementation for getting current status
	return &types.Status{}, nil
}

func DisplayStatus(status *types.Status) {
	// Implementation for displaying status in terminal
	fmt.Printf("PID: %d\n", status.PID)
	fmt.Printf("CPU Usage: %.2f%%\n", status.CPU.Usage)
	fmt.Printf("Memory Usage: %d MB\n", status.Memory.RSS/1024/1024)
	fmt.Printf("Event Loop Lag: %.2fms\n", status.EventLoop.Lag)
}