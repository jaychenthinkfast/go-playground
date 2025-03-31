package runner

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-playground/pkg/sandbox"
)

// Supported Go versions
const (
	Go124 = "go1.24"
	Go123 = "go1.23"
	Go122 = "go1.22"
)

// Version information
var Versions = map[string]string{
	Go124: "Go 1.24 - Released February 2024",
	Go123: "Go 1.23 - Released August 2023",
	Go122: "Go 1.22 - Released February 2023",
}

// Fixed time for deterministic output
var FixedTime = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)

// Run executes Go code in a sandbox
func Run(ctx context.Context, s *sandbox.Sandbox, code string, version string) (string, error) {
	// Validate version
	if !IsValidVersion(version) {
		return "", errors.New("unsupported Go version")
	}

	// Inject time package override for deterministic output if needed
	code = injectTimeOverride(code)

	// Check for security issues in the code
	if containsRestrictedOperations(code) {
		return "", errors.New("security sensitive operations are not allowed in the playground")
	}

	// Run the code
	return s.CompileAndRun(ctx, code, version)
}

// Format formats Go code
func Format(code string) (string, error) {
	s := sandbox.NewSandbox()
	return s.FormatCode(code)
}

// IsValidVersion checks if the provided Go version is supported
func IsValidVersion(version string) bool {
	_, ok := Versions[version]
	return ok
}

// injectTimeOverride modifies the code to use a fixed time
func injectTimeOverride(code string) string {
	// This is a simplified version. A real implementation would
	// use AST parsing and modification to properly inject the time override.

	// Only inject the override if the code explicitly imports "time" package
	// and not just any package that contains the string "time" (like "runtime")
	timeImportFound := false

	// Check for explicit time import
	if strings.Contains(code, "import") {
		// Check for standalone time import
		if strings.Contains(code, "import \"time\"") {
			timeImportFound = true
		}

		// Check for time import in a block
		if strings.Contains(code, "import (") &&
			(strings.Contains(code, "\"time\"") || strings.Contains(code, "\ttime")) {
			timeImportFound = true
		}
	}

	// Only proceed if time package is explicitly imported
	if !timeImportFound {
		return code
	}

	// Simplified - in real implementation, this would be more sophisticated
	override := `
// Playground time override
var _ = func() interface{} {
	// Set playground fixed time
	var timeNow = func() time.Time { return time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC) }
	return nil
}()
`
	// Find the position to insert the override (after imports)
	parts := strings.SplitN(code, "import", 2)
	if len(parts) == 2 {
		// Find the closing parenthesis if it's a block import
		if strings.Contains(parts[1], ")") {
			importParts := strings.SplitN(parts[1], ")", 2)
			if len(importParts) == 2 {
				return parts[0] + "import" + importParts[0] + ")" + override + importParts[1]
			}
		} else {
			// Handle single-line import
			newLineParts := strings.SplitN(parts[1], "\n", 2)
			if len(newLineParts) == 2 {
				return parts[0] + "import" + newLineParts[0] + "\n" + override + newLineParts[1]
			}
		}
	}

	return code
}

// containsRestrictedOperations checks for security-sensitive operations
func containsRestrictedOperations(code string) bool {
	restrictedOperations := []string{
		"os.Remove",
		"os.RemoveAll",
		"syscall.Exec",
		"syscall.ForkExec",
	}

	for _, op := range restrictedOperations {
		if strings.Contains(code, op) {
			return true
		}
	}

	return false
}
