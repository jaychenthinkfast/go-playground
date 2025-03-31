package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Sandbox represents a secure environment for running Go code
type Sandbox struct {
	TempDir    string
	MaxMemory  int64
	MaxCPUTime time.Duration
}

// NewSandbox creates a new sandbox with default limitations
func NewSandbox() *Sandbox {
	return &Sandbox{
		TempDir:    os.TempDir(),
		MaxMemory:  50 * 1024 * 1024, // 50MB
		MaxCPUTime: 5 * time.Second,
	}
}

// WithTimeout creates a context with a timeout
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// CompileAndRun compiles and runs Go code within the sandbox
func (s *Sandbox) CompileAndRun(ctx context.Context, code string, version string) (string, error) {
	// Create a temporary directory for the code
	dir, err := os.MkdirTemp(s.TempDir, "goplayground-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	// Initialize go.mod file for module support
	initCmd := exec.Command("go", "mod", "init", "playground")
	initCmd.Dir = dir
	if err := initCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Write the code to a file
	filename := filepath.Join(dir, "main.go")
	err = os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	// Fetch required dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = dir
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		return string(tidyOutput), fmt.Errorf("failed to fetch dependencies: %w", err)
	}

	// Select the Go binary and arguments based on version
	goBinary := "go"
	goArgs := []string{"run"}

	// Log version information - in a production environment we would use
	// different Go installations or containers for different versions
	// For this implementation, we'll just log the requested version
	versionInfo := fmt.Sprintf("Requested Go version: %s\n", version)

	// Add the filename to the arguments
	goArgs = append(goArgs, filename)

	// Run the code
	cmd := exec.CommandContext(ctx, goBinary, goArgs...)
	cmd.Dir = dir

	// Set resource limits in a real implementation
	// This is simplified for this example

	// Get real Go version
	var realVersionInfo bytes.Buffer
	versionCmd := exec.Command("go", "version")
	versionCmd.Stdout = &realVersionInfo
	versionCmd.Run()

	// Capture output
	output, err := cmd.CombinedOutput()

	// Format the full output with version info
	fullOutput := fmt.Sprintf("%sActual Go version: %s\n\n%s",
		versionInfo, strings.TrimSpace(realVersionInfo.String()), string(output))

	if err != nil {
		// Return both the error output and the error itself
		// This helps debug compilation issues
		return fullOutput, fmt.Errorf("%s: %w", string(output), err)
	}

	return fullOutput, nil
}

// FormatCode formats Go code using gofmt
func (s *Sandbox) FormatCode(code string) (string, error) {
	// Create a temporary directory for the code
	dir, err := os.MkdirTemp(s.TempDir, "goplayground-format-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	// Write the code to a file
	filename := filepath.Join(dir, "main.go")
	err = os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	// Format the code
	cmd := exec.Command("gofmt", filename)
	formattedCode, err := cmd.Output()
	if err != nil {
		return code, err
	}

	return string(formattedCode), nil
}
