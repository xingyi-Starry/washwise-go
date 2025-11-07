package util

import (
	"testing"
)

func TestLogger(t *testing.T) {
	InitLogger(LogConfig{
		Level: "trace",
		Dir:   "logs_test",
	})
}
