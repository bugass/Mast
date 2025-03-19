package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Progress struct {
	Total     int64
	Current   int64
	StartTime time.Time
	Speed     float64
}

type ProgressTracker struct {
	progress *Progress
	writer   io.Writer
}

func NewProgress(total int64) *ProgressTracker {
	return &ProgressTracker{
		progress: &Progress{
			Total:     total,
			StartTime: time.Now(),
		},
		writer: os.Stderr,
	}
}

func (pt *ProgressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.progress.Current += int64(n)
	pt.progress.Speed = float64(pt.progress.Current) / time.Since(pt.progress.StartTime).Seconds()
	pt.print()
	return n, nil
}

func (pt *ProgressTracker) Start() {
	pt.print()
}

func (pt *ProgressTracker) Stop() {
	fmt.Fprintf(pt.writer, "\n")
}

func (pt *ProgressTracker) print() {
	if pt.progress.Total == 0 {
		return
	}

	percent := float64(pt.progress.Current) / float64(pt.progress.Total) * 100
	speed := pt.progress.Speed / 1024 / 1024 // Convert to MB/s

	// Clear the current line
	fmt.Fprintf(pt.writer, "\r%s", strings.Repeat(" ", 80))

	// Print progress
	fmt.Fprintf(pt.writer, "\rProgress: %.1f%% (%.1f MB/s)", percent, speed)
}
