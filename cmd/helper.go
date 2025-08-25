package cmd

import (
	"fmt"
	"strings"

	"github.com/nexaa-cloud/nexaa-cli/api"
)

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

func envsToApi(environmentVariables []string, secret bool, state api.State) []api.EnvironmentVariableInput {
	var envs []api.EnvironmentVariableInput
	for _, env := range environmentVariables {
		parts := strings.Split(env, "=")
		name := parts[0]

		var value string
		if len(parts) != 1 {
			value = parts[1]
		}
		envs = append(envs, api.EnvironmentVariableInput{
			Name:   name,
			Value:  value,
			Secret: secret,
			State:  state,
		})
	}

	return envs
}
