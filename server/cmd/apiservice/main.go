package main

import (
	"github.com/asif10388/synctab/apiserver/controller"
)

func main() {
	mainRouter := controller.MakeController()
	mainRouter.Run(":5000")
}
