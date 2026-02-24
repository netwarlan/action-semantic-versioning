package action

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// SetOutput writes a GitHub Actions output variable.
func SetOutput(name, value string) error {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		// Not running in GitHub Actions — just print for debugging.
		fmt.Printf("::set-output name=%s::%s\n", name, value)
		return nil
	}

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open GITHUB_OUTPUT: %w", err)
	}
	defer func() { _ = f.Close() }()

	if strings.Contains(value, "\n") {
		delimiter := randomDelimiter()
		_, err = fmt.Fprintf(f, "%s<<%s\n%s\n%s\n", name, delimiter, value, delimiter)
	} else {
		_, err = fmt.Fprintf(f, "%s=%s\n", name, value)
	}

	return err
}

func randomDelimiter() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return "EOF_" + hex.EncodeToString(b)
}
