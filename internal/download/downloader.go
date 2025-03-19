package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bugass/mast/internal/download/progress"
)

const (
	defaultChunkSize  = 10 * 1024 * 1024
	defaultMaxRetries = 3
	defaultRetryDelay = 2 * time.Second
)

type DownloadConfig struct {
	ChunkSize     int64
	MaxRetries    int
	RetryDelay    time.Duration
	ResumeEnabled bool
}

type DownloadTask struct {
	URL         string
	Destination string
	Headers     map[string]string
	Cookies     []string
}

type Downloader struct {
	config     DownloadConfig
	httpClient *http.Client
}

func NewDownloader(config DownloadConfig) *Downloader {
	if config.ChunkSize == 0 {
		config.ChunkSize = defaultChunkSize
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}

	return &Downloader{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (d *Downloader) Download(task *DownloadTask) error {
	if task.CanResume() && d.config.ResumeEnabled {
		return d.resumeDownload(task)
	}
	return d.startNewDownload(task)
}

func (d *Downloader) startNewDownload(task *DownloadTask) error {
	req, err := http.NewRequest("GET", task.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mast/1.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	for key, value := range task.Headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range task.Cookies {
		parts := strings.Split(cookie, "=")
		if len(parts) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/") &&
		!strings.Contains(contentType, "application/pdf") &&
		!strings.Contains(contentType, "application/octet-stream") &&
		!strings.Contains(contentType, "application/x-msdownload") &&
		!strings.Contains(contentType, "application/x-download") &&
		!strings.Contains(contentType, "application/zip") &&
		!strings.Contains(contentType, "application/x-zip") &&
		!strings.Contains(contentType, "application/x-zip-compressed") {
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), "login") || strings.Contains(string(body), "authentication") {
			return fmt.Errorf("authentication failed: server returned login page instead of file")
		}
		return fmt.Errorf("unexpected content type: %s. Server might be requiring authentication", contentType)
	}

	destDir := filepath.Dir(task.Destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	file, err := os.Create(task.Destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer file.Close()

	progressTracker := progress.NewProgress(resp.ContentLength)
	progressTracker.Start()

	_, err = io.Copy(file, io.TeeReader(resp.Body, progressTracker))
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	progressTracker.Stop()
	return nil
}

func (d *Downloader) resumeDownload(task *DownloadTask) error {
	fileInfo, err := os.Stat(task.Destination)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	req, err := http.NewRequest("GET", task.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mast/1.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileInfo.Size()))

	for key, value := range task.Headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range task.Cookies {
		parts := strings.Split(cookie, "=")
		if len(parts) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server does not support resume: %s", resp.Status)
	}

	file, err := os.OpenFile(task.Destination, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	progressTracker := progress.NewProgress(resp.ContentLength + fileInfo.Size())
	progressTracker.Start()

	_, err = io.Copy(file, io.TeeReader(resp.Body, progressTracker))
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	progressTracker.Stop()
	return nil
}

func (t *DownloadTask) CanResume() bool {
	_, err := os.Stat(t.Destination)
	return err == nil
}
