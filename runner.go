package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const terminalTitlePrefix = "Mado-Tray"

func madoTerminalTitle(scriptID string) string {
	id := strings.TrimSpace(scriptID)
	if id == "" {
		return terminalTitlePrefix
	}
	return fmt.Sprintf("%s: %s", terminalTitlePrefix, id)
}

func IsMadoTerminalOpen(title string) (bool, error) {
	appleScript := fmt.Sprintf(`set found to false
if application "Terminal" is running then
	tell application "Terminal"
		repeat with terminalWindow in windows
			repeat with terminalTab in tabs of terminalWindow
				try
					if custom title of terminalTab is equal to "%s" then
						set found to true
						exit repeat
					end if
				end try
			end repeat
			if found then exit repeat
		end repeat
	end tell
end if
return found`, escapeAppleScriptString(title))

	output, err := exec.Command("osascript", "-e", appleScript).Output()
	if err != nil {
		return false, fmt.Errorf("consultando Terminal.app: %w", err)
	}

	return strings.TrimSpace(string(output)) == "true", nil
}

func GetRunningMadoScriptIDs() (map[string]bool, error) {
	prefix := terminalTitlePrefix + ": "
	appleScript := fmt.Sprintf(`set found to ""
if application "Terminal" is running then
	tell application "Terminal"
		repeat with terminalWindow in windows
			repeat with terminalTab in tabs of terminalWindow
				try
					set tabTitle to custom title of terminalTab
					if tabTitle starts with "%s" and busy of terminalTab is true then
						if found is not "" then set found to found & linefeed
						set found to found & tabTitle
					end if
				end try
			end repeat
		end repeat
	end tell
end if
return found`, escapeAppleScriptString(prefix))

	output, err := exec.Command("osascript", "-e", appleScript).Output()
	if err != nil {
		return nil, fmt.Errorf("querying Terminal.app: %w", err)
	}

	running := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		title := strings.TrimSpace(line)
		if !strings.HasPrefix(title, prefix) {
			continue
		}
		id := strings.TrimSpace(title[len(prefix):])
		if id != "" {
			running[id] = true
		}
	}

	return running, nil
}

func RunInVisibleTerminal(scriptPath, title string) error {
	if strings.TrimSpace(scriptPath) == "" {
		return fmt.Errorf("la ruta del script está vacía")
	}

	if strings.TrimSpace(title) == "" {
		title = terminalTitlePrefix
	}

	command, err := buildTerminalCommand(scriptPath)
	if err != nil {
		return err
	}
	appleScript := fmt.Sprintf(`tell application "Terminal"
	set madoTab to do script "%s"
	set custom title of madoTab to "%s"
end tell`, escapeAppleScriptString(command), escapeAppleScriptString(title))

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
