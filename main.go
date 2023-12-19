package main

import (
	"bub/src"
	"fmt"
	"os"
)

func main() {

	session := src.New()

	err := session.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
