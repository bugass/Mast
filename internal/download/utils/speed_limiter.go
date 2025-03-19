package utils

import (
	"io"
	"time"
)

type SpeedLimitReader struct {
	reader io.Reader
	limit  int64
}

func NewSpeedLimitReader(reader io.Reader, limit int64) *SpeedLimitReader {
	return &SpeedLimitReader{
		reader: reader,
		limit:  limit,
	}
}

func (r *SpeedLimitReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		time.Sleep(time.Duration(float64(n) / float64(r.limit) * float64(time.Second)))
	}
	return
}
