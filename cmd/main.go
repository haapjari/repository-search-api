package main

import (
	"github.com/haapjari/glass/pkg/router"
	"log"
)

func main() {
	if err := router.SetupRouter(); err != nil {
		log.Fatal(err.Error())
	}
}
