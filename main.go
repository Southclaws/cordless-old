package main

import (
	"fmt"

	"github.com/Southclaws/cordless/core"
)

func main() {
	app, err := core.Initialise()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
