package main

import (
	"fmt"

	"github.com/asif10388/synctab/apiserver"
)

func main() {

	input := apiserver.APIServerInput{
		DeploymentType: string("dev"),
	}

	myAPIServer := apiserver.NewApiServer(input)

	err := myAPIServer.Init()
	if err != nil {
		fmt.Errorf("failed to initialize API service")
	}

	myAPIServer.Start()

}
