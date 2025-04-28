package main

import (
	"log"

	"gitlab.com/tilaa/tilaa-cli/cmd"
)

func main() {
	log.SetFlags(0)

	cmd.Execute()
}

//go:generate go run github.com/Khan/genqlient genqlient.yaml
