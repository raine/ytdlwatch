package main

import (
	"os"
	"strings"
)

func isDevelopment() bool {
	tmpDir, ok := os.LookupEnv("TMPDIR")
	return ok && strings.HasPrefix(os.Args[0], tmpDir)
}
