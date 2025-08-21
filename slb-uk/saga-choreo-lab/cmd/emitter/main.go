package main

import (
	"log"
	"example.com/saga-choreo-lab/pkg/common"
)

func main() {
	if err := common.RunEmitter(); err != nil {
		log.Fatal(err)
	}
}
