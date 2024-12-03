package main

import (
	"log"

	"gitlab.com/Tilaa/tilaa-cli/cmd"
)

func main() {
	log.SetFlags(0)

	cmd.Execute()
}
