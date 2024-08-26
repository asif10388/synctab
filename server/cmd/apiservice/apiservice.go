package main

import (
	"fmt"
	"os"

	"github.com/asif10388/synctab/apiserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	input := apiserver.APIServerInput{
		DeploymentType: string("dev"),
	}

	myAPIServer, err := apiserver.NewApiServer(input)
	if err != nil {
		fmt.Println("failed to allocate API service")
	}

	err = myAPIServer.Init()
	if err != nil {
		fmt.Println("failed to initialize API service")
	}

	myAPIServer.Start()

}
