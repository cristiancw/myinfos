package main

import (
	"fmt"
	"time"

	"github.com/cristiancw/myinfos/info"
)

func main() {
	machineChan := make(chan info.Machine)
	go info.LoadMachine(time.Now(), machineChan)
	defer close(machineChan)

	for machine := range machineChan {
		fmt.Printf("1-%v\n", machine)
	}
}
