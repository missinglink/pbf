package lib

import (
	"fmt"
	"os"
)

// EnsureDirectoryExists - otherwise fatally fail
func EnsureDirectoryExists(path string, label string) os.FileInfo {

	// stat destination
	info, err := os.Stat(path)

	// path not found
	if err != nil {
		fmt.Println(label, "path does not exist")
		os.Exit(1)
	}

	// not a directory
	if !info.IsDir() {
		fmt.Println(label, " path not a directory")
		os.Exit(1)
	}

	return info
}
