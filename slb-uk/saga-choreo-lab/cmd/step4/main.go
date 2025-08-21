package main

import (
	"log"
	"example.com/saga-choreo-lab/pkg/common"
)

func main() {
	if err := common.RunStepService(); err != nil {
		log.Fatal(err)
	}
}
