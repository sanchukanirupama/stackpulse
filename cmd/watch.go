package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"stackpulse/internal/monitor"
	"stackpulse/internal/config"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start monitoring a Node.js service",
	Long: `Monitor a Node.js service for memory leaks, CPU spikes, and performance issues.
	
Examples:
  stackpulse watch --pid 1234 --cpu-threshold 80 --heap-limit 200MB
  stackpulse watch --port 3000 --polling-ms 100`,
	RunE: runWatch,
}

var (
	host          string
	port          int
	pid           int
	heapLimit     string
	cpuThreshold  float64
	pollingMs     int
	inspectPort   int
)

func init() {
	rootCmd.AddCommand(watchCmd)
	
	watchCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Host to monitor")
	watchCmd.Flags().IntVar(&port, "port", 0, "Port to monitor")
	watchCmd.Flags().IntVar(&pid, "pid", 0, "Process ID to monitor")
	watchCmd.Flags().StringVar(&heapLimit, "heap-limit", "150MB", "Heap memory limit threshold")
	watchCmd.Flags().Float64Var(&cpuThreshold, "cpu-threshold", 70.0, "CPU usage threshold percentage")
	watchCmd.Flags().IntVar(&pollingMs, "polling-ms", 100, "Polling interval in milliseconds")
	watchCmd.Flags().IntVar(&inspectPort, "inspect-port", 9229, "V8 inspector port")
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfg := &config.ServiceConfig{
		Host:            host,
		Port:            port,
		PID:             pid,
		InspectPort:     inspectPort,
		HeapLimit:       heapLimit,
		CPUThreshold:    cpuThreshold,
		PollingInterval: time.Duration(pollingMs) * time.Millisecond,
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	monitor := monitor.New(cfg)
	
	go func() {
		<-sigChan
		fmt.Println("\nShutting down gracefully...")
		cancel()
	}()

	return monitor.Start(ctx)
}