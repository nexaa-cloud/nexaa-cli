package main

import (
	"log"

	"flag"
	"fmt"
	"os"

	"gitlab.com/tilaa/tilaa-cli/cmd"
	"gitlab.com/tilaa/tilaa-cli/config"
)

func main() {
	env := flag.String("env", config.GetEnvironment(), "Environment (dev, prod)")
	flag.Parse()

	if err := config.Initialize(*env); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	log.SetFlags(0)

	cmd.Execute()
}

//go:generate go run github.com/Khan/genqlient genqlient.yaml
