package display

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"stackpulse/internal/types"
)

type Dashboard struct {
	lastUpdate time.Time
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (d *Dashboard) Update(status *types.Status) {
	d.clearScreen()
	d.displayHeader()
	d.displayMetrics(status)
	d.displayAlerts(status.Alerts)
	d.lastUpdate = time.Now()
}

func (d *Dashboard) clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (d *Dashboard) displayHeader() {
	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	headerColor.Println("║                            STACKPULSE DASHBOARD                              ║")
	headerColor.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Printf("Last Update: %s\n\n", d.lastUpdate.Format("15:04:05.000"))
}

func (d *Dashboard) displayMetrics(status *types.Status) {
	// Service info
	serviceColor := color.New(color.FgGreen, color.Bold)
	serviceColor.Printf("🔍 Monitoring PID: %d\n\n", status.PID)

	// Create table for metrics
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Current", "Status", "Threshold"})
	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	// CPU metrics
	cpuStatus := "✅ Normal"
	cpuColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if status.CPU.Usage > 70 {
		cpuStatus = "⚠️  High"
		cpuColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if status.CPU.Usage > 90 {
		cpuStatus = "🚨 Critical"
		cpuColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"CPU Usage",
		fmt.Sprintf("%.2f%%", status.CPU.Usage),
		cpuStatus,
		"< 70%",
	}, []tablewriter.Colors{{}, cpuColor, cpuColor, {}})

	// Memory metrics
	memoryMB := float64(status.Memory.RSS) / 1024 / 1024
	memoryStatus := "✅ Normal"
	memoryColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if memoryMB > 100 {
		memoryStatus = "⚠️  High"
		memoryColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if memoryMB > 200 {
		memoryStatus = "🚨 Critical"
		memoryColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"Memory (RSS)",
		fmt.Sprintf("%.1f MB", memoryMB),
		memoryStatus,
		"< 150 MB",
	}, []tablewriter.Colors{{}, memoryColor, memoryColor, {}})

	// Heap metrics
	if status.Memory.HeapTotal > 0 {
		heapUsedMB := float64(status.Memory.HeapUsed) / 1024 / 1024
		heapTotalMB := float64(status.Memory.HeapTotal) / 1024 / 1024
		heapUsage := (float64(status.Memory.HeapUsed) / float64(status.Memory.HeapTotal)) * 100

		heapStatus := "✅ Normal"
		heapColor := tablewriter.Colors{tablewriter.FgGreenColor}
		if heapUsage > 80 {
			heapStatus = "⚠️  High"
			heapColor = tablewriter.Colors{tablewriter.FgYellowColor}
		}
		if heapUsage > 95 {
			heapStatus = "🚨 Critical"
			heapColor = tablewriter.Colors{tablewriter.FgRedColor}
		}

		table.Rich([]string{
			"Heap Usage",
			fmt.Sprintf("%.1f/%.1f MB (%.1f%%)", heapUsedMB, heapTotalMB, heapUsage),
			heapStatus,
			"< 80%",
		}, []tablewriter.Colors{{}, heapColor, heapColor, {}})
	}

	// Event loop lag
	lagStatus := "✅ Normal"
	lagColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if status.EventLoop.Lag > 5 {
		lagStatus = "⚠️  High"
		lagColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if status.EventLoop.Lag > 20 {
		lagStatus = "🚨 Critical"
		lagColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"Event Loop Lag",
		fmt.Sprintf("%.2f ms", status.EventLoop.Lag),
		lagStatus,
		"< 5 ms",
	}, []tablewriter.Colors{{}, lagColor, lagColor, {}})

	// Event loop utilization
	utilizationStatus := "✅ Normal"
	utilizationColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if status.EventLoop.Utilization > 70 {
		utilizationStatus = "⚠️  High"
		utilizationColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if status.EventLoop.Utilization > 90 {
		utilizationStatus = "🚨 Critical"
		utilizationColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"Event Loop Util",
		fmt.Sprintf("%.1f%%", status.EventLoop.Utilization),
		utilizationStatus,
		"< 70%",
	}, []tablewriter.Colors{{}, utilizationColor, utilizationColor, {}})

	// GC metrics
	gcStatus := "✅ Normal"
	gcColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if status.GC.Duration > 10 {
		gcStatus = "⚠️  High"
		gcColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if status.GC.Duration > 50 {
		gcStatus = "🚨 Critical"
		gcColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"GC Duration",
		fmt.Sprintf("%.2f ms (%s)", status.GC.Duration, status.GC.Type),
		gcStatus,
		"< 10 ms",
	}, []tablewriter.Colors{{}, gcColor, gcColor, {}})

	// Handle metrics
	handleStatus := "✅ Normal"
	handleColor := tablewriter.Colors{tablewriter.FgGreenColor}
	if status.Handles.Active > 50 {
		handleStatus = "⚠️  High"
		handleColor = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if status.Handles.Active > 100 {
		handleStatus = "🚨 Critical"
		handleColor = tablewriter.Colors{tablewriter.FgRedColor}
	}

	table.Rich([]string{
		"Active Handles",
		fmt.Sprintf("%d (T:%d, S:%d)", status.Handles.Active, status.Handles.Timers, status.Handles.TCPSockets),
		handleStatus,
		"< 50",
	}, []tablewriter.Colors{{}, handleColor, handleColor, {}})

	table.Render()
	fmt.Println()

	// Display additional metrics in a second table
	d.displayAdvancedMetrics(status)
}

func (d *Dashboard) displayAdvancedMetrics(status *types.Status) {
	advancedColor := color.New(color.FgMagenta, color.Bold)
	advancedColor.Println("📊 Advanced Node.js Metrics:")

	// Create advanced metrics table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Current", "Details"})
	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgMagentaColor},
	)

	// Event loop statistics
	table.Append([]string{
		"Event Loop Stats",
		fmt.Sprintf("Avg: %.2fms", status.EventLoop.Mean),
		fmt.Sprintf("Min: %.2f, Max: %.2f, P95: %.2f", 
			status.EventLoop.Min, status.EventLoop.Max, status.EventLoop.P95),
	})

	// Thread pool
	table.Append([]string{
		"Thread Pool",
		fmt.Sprintf("Active: %d/%d", status.ThreadPool.ActiveCount, status.ThreadPool.PoolSize),
		fmt.Sprintf("Queue: %d, Pending: %d", 
			status.ThreadPool.QueueSize, status.ThreadPool.PendingCount),
	})

	// GC statistics
	table.Append([]string{
		"Garbage Collection",
		fmt.Sprintf("Collections: %d", status.GC.Collections),
		fmt.Sprintf("Total: %d (%.2fms), Reason: %s", 
			status.GC.CollectionsTotal, status.GC.DurationTotal, status.GC.Reason),
	})

	// V8 heap spaces
	if len(status.V8.HeapSpaceUsed) > 0 {
		var heapDetails []string
		for space, used := range status.V8.HeapSpaceUsed {
			heapDetails = append(heapDetails, fmt.Sprintf("%s: %.1fMB", 
				space, float64(used)/1024/1024))
		}
		table.Append([]string{
			"V8 Heap Spaces",
			fmt.Sprintf("%d spaces", len(status.V8.HeapSpaceUsed)),
			strings.Join(heapDetails, ", "),
		})
	}

	// Memory details
	table.Append([]string{
		"Memory Details",
		fmt.Sprintf("Malloc: %.1fMB", float64(status.V8.MallocedMemory)/1024/1024),
		fmt.Sprintf("Peak: %.1fMB, External: %.1fMB", 
			float64(status.V8.PeakMallocedMemory)/1024/1024,
			float64(status.Memory.External)/1024/1024),
	})

	table.Render()
	fmt.Println()
}

func (d *Dashboard) displayAlerts(alerts []types.Alert) {
	if len(alerts) == 0 {
		successColor := color.New(color.FgGreen)
		successColor.Println("✅ No active alerts")
		return
	}

	alertColor := color.New(color.FgRed, color.Bold)
	alertColor.Printf("🚨 Active Alerts (%d):\n", len(alerts))
	
	for i, alert := range alerts {
		fmt.Printf("  %d. [%s] %s (%.2f > %.2f)\n", 
			i+1, 
			string(alert.Severity), 
			alert.Message, 
			alert.Value, 
			alert.Threshold)
	}
	fmt.Println()
}