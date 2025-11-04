package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"time"
)

func Logger(level string, message string, fields ...map[string]any) {
	data := map[string]any{
		"level":     level,
		"message":   message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if len(fields) > 0 {
		for _, f := range fields {
			maps.Copy(data, f)
		}
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

func Fatal(message string, fields ...map[string]any) {
	Logger("error", message, fields...)
	os.Exit(1)
}
