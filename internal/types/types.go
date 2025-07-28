package types

import (
	"time"
)

// AlertType represents the type of alert
type AlertType string

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	AlertTypeCPU       AlertType = "cpu"
	AlertTypeMemory    AlertType = "memory"
	AlertTypeEventLoop AlertType = "eventloop"
	AlertTypeHeap      AlertType = "heap"

	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents a monitoring alert
type Alert struct {
	Type      AlertType     `json:"type"`
	Severity  AlertSeverity `json:"severity"`
	Message   string        `json:"message"`
	Value     float64       `json:"value"`
	Threshold float64       `json:"threshold"`
	Timestamp time.Time     `json:"timestamp"`
}

// CPUMetrics represents CPU usage metrics
type CPUMetrics struct {
	Usage      float64   `json:"usage"`
	UserTime   float64   `json:"userTime"`
	SystemTime float64   `json:"systemTime"`
	Timestamp  time.Time `json:"timestamp"`
}

// MemoryMetrics represents memory usage metrics
type MemoryMetrics struct {
	RSS        uint64    `json:"rss"`
	VMS        uint64    `json:"vms"`
	HeapTotal  uint64    `json:"heapTotal"`
	HeapUsed   uint64    `json:"heapUsed"`
	External   uint64    `json:"external"`
	Timestamp  time.Time `json:"timestamp"`
}

// EventLoopMetrics represents event loop performance metrics
type EventLoopMetrics struct {
	Lag         float64   `json:"lag"`
	Mean        float64   `json:"mean"`
	Max         float64   `json:"max"`
	P95         float64   `json:"p95"`
	Min         float64   `json:"min"`
	Utilization float64   `json:"utilization"`
	Timestamp   time.Time `json:"timestamp"`
}

// ThreadPoolMetrics represents thread pool metrics
type ThreadPoolMetrics struct {
	QueueSize    int       `json:"queueSize"`
	PoolSize     int       `json:"poolSize"`
	ActiveCount  int       `json:"activeCount"`
	PendingCount int       `json:"pendingCount"`
	Timestamp    time.Time `json:"timestamp"`
}

// GCMetrics represents garbage collection metrics
type GCMetrics struct {
	Collections      int       `json:"collections"`
	Duration         float64   `json:"duration"`
	HeapSizeBefore   uint64    `json:"heapSizeBefore"`
	HeapSizeAfter    uint64    `json:"heapSizeAfter"`
	Type             string    `json:"type"`
	Reason           string    `json:"reason"`
	CollectionsTotal int       `json:"collectionsTotal"`
	DurationTotal    float64   `json:"durationTotal"`
	Timestamp        time.Time `json:"timestamp"`
}

// HandleMetrics represents handle usage metrics
type HandleMetrics struct {
	Active    int       `json:"active"`
	Refs      int       `json:"refs"`
	Timers    int       `json:"timers"`
	TCPSockets int      `json:"tcpSockets"`
	UDPSockets int      `json:"udpSockets"`
	Files     int       `json:"files"`
	Timestamp time.Time `json:"timestamp"`
}

// V8Metrics represents V8 engine specific metrics
type V8Metrics struct {
	HeapSpaceUsed      map[string]uint64 `json:"heapSpaceUsed"`
	HeapSpaceSize      map[string]uint64 `json:"heapSpaceSize"`
	HeapSpaceAvailable map[string]uint64 `json:"heapSpaceAvailable"`
	MallocedMemory     uint64            `json:"mallocedMemory"`
	PeakMallocedMemory uint64            `json:"peakMallocedMemory"`
	Timestamp          time.Time         `json:"timestamp"`
}

// Status represents the current monitoring status
type Status struct {
	PID         int               `json:"pid"`
	CPU         CPUMetrics        `json:"cpu"`
	Memory      MemoryMetrics     `json:"memory"`
	EventLoop   EventLoopMetrics  `json:"eventLoop"`
	ThreadPool  ThreadPoolMetrics `json:"threadPool"`
	GC          GCMetrics         `json:"gc"`
	Handles     HandleMetrics     `json:"handles"`
	V8          V8Metrics         `json:"v8"`
	Timestamp   time.Time         `json:"timestamp"`
	Alerts      []Alert           `json:"alerts"`
}