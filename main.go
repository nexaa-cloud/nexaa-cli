package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nexaa-cloud/nexaa-cli/cmd"
	"github.com/nexaa-cloud/nexaa-cli/config"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	config.Initialize()

	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	log.SetFlags(0)

	cmd.Execute()
}

//go:generate go run -mod=mod github.com/Khan/genqlient genqlient.yaml
