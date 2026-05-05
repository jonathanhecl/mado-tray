package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunInVisibleTerminal(scriptPath string) error {
	if strings.TrimSpace(scriptPath) == "" {
		return fmt.Errorf("la ruta del script está vacía")
	}

	command := buildTerminalCommand(scriptPath)
	appleScript := fmt.Sprintf(`tell application "Terminal"
	do script "%s"
	activate
end tell`, escapeAppleScriptString(command))

	cmd := exec.Command("osascript", "-e", appleScript)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("abriendo Terminal.app: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func buildTerminalCommand(scriptPath string) string {
	cleanPath := filepath.Clean(scriptPath)
	dir := filepath.Dir(cleanPath)

	if _, err := os.Stat(cleanPath); err == nil {
		return fmt.Sprintf("cd %s && %s", shellQuote(dir), shellQuote(cleanPath))
	}

	return shellQuote(cleanPath)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func escapeAppleScriptString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	return strings.ReplaceAll(value, `"`, `\"`)
}
