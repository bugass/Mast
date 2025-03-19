package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mast",
	Short: "A fast and reliable file downloader",
	Long:  `Mast is a command-line tool for downloading files with support for resumable downloads, custom headers, and cookies.`,
}

func Execute() error {
	rootCmd.AddCommand(versionCmd)
	return rootCmd.Execute()
}
