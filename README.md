# Tilaa CLI

The Tilaa CLI is a tool to manage cloud resources on the Tilaa Serverless
Platform, which allows you to deploy containers very quickly.

## Setup your environment

Make sure you have golang (>=1.23.0) running (either locally or in a container).

## Building the binary

Building the binary can simply be done by `go build -o tilaa .`

## Running development

If you're developing and want to talk to staging instead of production, use the environment variable TILAA_ENV=dev to talk to the dev environment

`go run . <args>`

This will ensure the code will talk to staging-auth and staging-graphql. Your
access token, refresh token, and ExpiresAt will be stored in JSON in
`./auth.json`. In the production version, this is stored at "~/.tilaa/auth.json"
which is better.

## Autocomplete
To enable shell completion, run:

For Bash:
    `source <(tilaa-cli completion bash)`

For Zsh:
    `source <(tilaa-cli completion zsh)`

Or to persist it, save the output to a file and source it in your shell config.