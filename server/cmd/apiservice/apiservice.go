package main

import (
	"fmt"

	"github.com/asif10388/synctab/apiserver"
)

func main() {
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
