package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bugass/mast/internal/download"
	"github.com/spf13/cobra"
)

var (
	dest     string
	location string
	cookies  []string
	headers  []string
	retries  int
	resume   bool
)

var downloadCmd = &cobra.Command{
	Use:   "download [URL]",
	Short: "Download a file from a URL",
	Long: `Download a file from a URL with optional flags for customization.
Example:
  mast download https://example.com/file.zip -f file.zip -l downloads/
  mast download https://example.com/file.zip --header "Authorization: Bearer token"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		if dest == "" {
			dest = filepath.Base(url)
		}

		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			parts := strings.Split(cookie, "=")
			if len(parts) == 2 {
				cookieMap[parts[0]] = parts[1]
			}
		}

		headerMap := make(map[string]string)
		for _, header := range headers {
			parts := strings.Split(header, ":")
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		var finalDest string
		if location != "" {
			if err := os.MkdirAll(location, 0755); err != nil {
				return fmt.Errorf("failed to create location directory: %v", err)
			}
			finalDest = filepath.Join(location, dest)
		} else {
			finalDest = dest
		}

		config := download.DownloadConfig{
			ChunkSize:     10 * 1024 * 1024,
			MaxRetries:    retries,
			RetryDelay:    2 * time.Second,
			ResumeEnabled: resume,
		}

		downloader := download.NewDownloader(config)
		task := &download.DownloadTask{
			URL:         url,
			Destination: finalDest,
			Headers:     headerMap,
			Cookies:     cookies,
		}

		if err := downloader.Download(task); err != nil {
			return fmt.Errorf("download failed: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&dest, "file", "f", "", "Destination filename (default: filename from URL)")
	downloadCmd.Flags().StringVarP(&location, "location", "l", "", "Location to save the file (optional)")
	downloadCmd.Flags().StringSliceVar(&cookies, "cookie", []string{}, "Cookies to send with the request (format: name=value)")
	downloadCmd.Flags().StringSliceVar(&headers, "header", []string{}, "Headers to send with the request (format: name:value)")
	downloadCmd.Flags().IntVar(&retries, "retries", 3, "Maximum number of retry attempts")
	downloadCmd.Flags().BoolVar(&resume, "resume", true, "Enable resumable downloads")
}
