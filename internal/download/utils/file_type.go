package utils

import (
	"path/filepath"
	"strings"
)

type FileType struct {
	Extension string
	MimeType  string
	Category  string
}

func DetectFileType(url, contentType string) *FileType {
	ft := &FileType{
		MimeType: contentType,
	}

	ext := strings.ToLower(filepath.Ext(url))
	if ext != "" {
		ft.Extension = ext[1:] 
	}

	switch {
	case strings.Contains(contentType, "pdf"):
		ft.Category = "document"
	case strings.Contains(contentType, "image"):
		ft.Category = "image"
	case strings.Contains(contentType, "video"):
		ft.Category = "video"
	case strings.Contains(contentType, "audio"):
		ft.Category = "audio"
	case strings.Contains(contentType, "text"):
		ft.Category = "text"
	case strings.Contains(contentType, "application/zip"):
		ft.Category = "archive"
	case strings.Contains(contentType, "application/x-rar-compressed"):
		ft.Category = "archive"
	case strings.Contains(contentType, "application/x-7z-compressed"):
		ft.Category = "archive"
	default:
		ft.Category = "unknown"
	}

	return ft
}
