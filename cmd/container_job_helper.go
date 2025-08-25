package cmd

import "fmt"

func commandApiToString(apiCommand []string) string {
	command := "["
	for i, cmd := range apiCommand {
		command += fmt.Sprintf(`%q`, cmd)

		if i < len(apiCommand)-1 {
			command += ","
		}
	}
	command += "]"
	return command
}

func enabledApiToString(enabled bool) string {
	if enabled {
		return "True"
	}
	return "False"
}
