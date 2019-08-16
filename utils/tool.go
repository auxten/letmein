package utils

import (
	"path/filepath"
	"runtime"
)

var FJ = filepath.Join

// GetProjectSrcDir gets the src code root.
func GetProjectSrcDir() string {
	_, testFile, _, _ := runtime.Caller(0)
	return FJ(filepath.Dir(testFile), "../")
}
