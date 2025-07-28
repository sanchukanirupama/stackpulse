package alerts

import (
	"fmt"
	"time"
	"stackpulse/internal/config"
	"stackpulse/internal/types"
)

type Manager struct {
	activeAlerts map[string]types.Alert
}

func NewManager() *Manager {
	return &Manager{
		activeAlerts: make(map[string]types.Alert),
	}
}

func (m *Manager) CheckThresholds(status *types.Status, cfg *config.ServiceConfig) []types.Alert {
	var alerts []types.Alert

	// Check CPU threshold
	if status.CPU.Usage > cfg.CPUThreshold {
		severity := types.SeverityWarning
		if status.CPU.Usage > 90 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      types.AlertTypeCPU,
			Severity:  severity,
			Message:   fmt.Sprintf("High CPU usage: %.2f%% (threshold: %.2f%%)", status.CPU.Usage, cfg.CPUThreshold),
			Value:     status.CPU.Usage,
			Threshold: cfg.CPUThreshold,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check memory threshold (simplified - would parse cfg.HeapLimit in production)
	memoryMB := float64(status.Memory.RSS) / 1024 / 1024
	if memoryMB > 150 { // Simplified threshold
		severity := types.SeverityWarning
		if memoryMB > 200 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      types.AlertTypeMemory,
			Severity:  severity,
			Message:   fmt.Sprintf("High memory usage: %.1f MB (threshold: 150 MB)", memoryMB),
			Value:     memoryMB,
			Threshold: 150,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check heap usage
	if status.Memory.HeapTotal > 0 {
		heapUsage := (float64(status.Memory.HeapUsed) / float64(status.Memory.HeapTotal)) * 100
		if heapUsage > 80 {
			severity := types.SeverityWarning
			if heapUsage > 95 {
				severity = types.SeverityCritical
			}
			
			alert := types.Alert{
				Type:      types.AlertTypeHeap,
				Severity:  severity,
				Message:   fmt.Sprintf("High heap usage: %.1f%% (threshold: 80%%)", heapUsage),
				Value:     heapUsage,
				Threshold: 80,
				Timestamp: time.Now(),
			}
			alerts = append(alerts, alert)
		}
	}

	// Check event loop lag
	if status.EventLoop.Lag > 5 {
		severity := types.SeverityWarning
		if status.EventLoop.Lag > 20 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      types.AlertTypeEventLoop,
			Severity:  severity,
			Message:   fmt.Sprintf("High event loop lag: %.2fms (threshold: 5ms)", status.EventLoop.Lag),
			Value:     status.EventLoop.Lag,
			Threshold: 5,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check event loop utilization
	if status.EventLoop.Utilization > 70 {
		severity := types.SeverityWarning
		if status.EventLoop.Utilization > 90 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      types.AlertTypeEventLoop,
			Severity:  severity,
			Message:   fmt.Sprintf("High event loop utilization: %.1f%% (threshold: 70%%)", status.EventLoop.Utilization),
			Value:     status.EventLoop.Utilization,
			Threshold: 70,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check GC duration
	if status.GC.Duration > 10 {
		severity := types.SeverityWarning
		if status.GC.Duration > 50 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      "gc",
			Severity:  severity,
			Message:   fmt.Sprintf("Long GC duration: %.2fms (threshold: 10ms)", status.GC.Duration),
			Value:     status.GC.Duration,
			Threshold: 10,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	// Check handle count
	if status.Handles.Active > 50 {
		severity := types.SeverityWarning
		if status.Handles.Active > 100 {
			severity = types.SeverityCritical
		}
		
		alert := types.Alert{
			Type:      "handles",
			Severity:  severity,
			Message:   fmt.Sprintf("High handle count: %d (threshold: 50)", status.Handles.Active),
			Value:     float64(status.Handles.Active),
			Threshold: 50,
			Timestamp: time.Now(),
		}
		alerts = append(alerts, alert)
	}

	return alerts
}