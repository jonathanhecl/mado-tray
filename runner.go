package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const terminalTitlePrefix = "Mado-Tray"

func RunInVisibleTerminal(scriptPath string) error {
	if strings.TrimSpace(scriptPath) == "" {
		return fmt.Errorf("la ruta del script está vacía")
	}

	command, err := buildTerminalCommand(scriptPath)
	if err != nil {
		return err
	}
	appleScript := fmt.Sprintf(`tell application "Terminal"
	set madoTab to do script "%s"
	set custom title of madoTab to "%s"
	activate
end tell`, escapeAppleScriptString(command), escapeAppleScriptString(terminalTitle(scriptPath)))

	cmd := exec.Command("osascript", "-e", appleScript)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("abriendo Terminal.app: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func CloseInactiveMadoTerminals() error {
	appleScript := fmt.Sprintf(`if application "Terminal" is running then
	tell application "Terminal"
		repeat with windowIndex from (count of windows) to 1 by -1
			set terminalWindow to window windowIndex
			repeat with tabIndex from (count of tabs of terminalWindow) to 1 by -1
				set terminalTab to tab tabIndex of terminalWindow
				try
					set tabTitle to custom title of terminalTab
					if tabTitle starts with "%s" and busy of terminalTab is false then
						close terminalTab
					end if
				end try
			end repeat
		end repeat
	end tell
end if`, escapeAppleScriptString(terminalTitlePrefix))

	cmd := exec.Command("osascript", "-e", appleScript)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cerrando terminales inactivas: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func buildTerminalCommand(scriptPath string) (string, error) {
	parts, err := splitShellWords(scriptPath)
	if err != nil {
		return "", fmt.Errorf("la ruta del script contiene comillas sin cerrar")
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("la ruta del script está vacía")
	}

	executable := filepath.Clean(parts[0])
	args := parts[1:]
	quotedCommand := make([]string, 0, len(parts))
	quotedCommand = append(quotedCommand, shellQuote(executable))
	for _, arg := range args {
		quotedCommand = append(quotedCommand, shellQuote(arg))
	}

	if _, err := os.Stat(executable); err == nil {
		dir := filepath.Dir(executable)
		return fmt.Sprintf("cd %s && %s", shellQuote(dir), strings.Join(quotedCommand, " ")), nil
	}

	return strings.Join(quotedCommand, " "), nil
}

func terminalTitle(scriptPath string) string {
	parts, err := splitShellWords(scriptPath)
	if err != nil || len(parts) == 0 {
		return terminalTitlePrefix
	}

	name := filepath.Base(filepath.Clean(parts[0]))
	if name == "." || name == string(filepath.Separator) {
		return terminalTitlePrefix
	}

	return fmt.Sprintf("%s: %s", terminalTitlePrefix, name)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func escapeAppleScriptString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	return strings.ReplaceAll(value, `"`, `\"`)
}

func splitShellWords(input string) ([]string, error) {
	var parts []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	for _, char := range strings.TrimSpace(input) {
		switch {
		case escaped:
			current.WriteRune(char)
			escaped = false
		case char == '\\' && !inSingle:
			escaped = true
		case char == '\'' && !inDouble:
			inSingle = !inSingle
		case char == '"' && !inSingle:
			inDouble = !inDouble
		case (char == ' ' || char == '\t') && !inSingle && !inDouble:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if escaped || inSingle || inDouble {
		return nil, fmt.Errorf("entrada inválida")
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts, nil
}
