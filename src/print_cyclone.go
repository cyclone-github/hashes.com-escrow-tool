package main

import (
	"fmt"
	"time"
)

func printCyclone() {
	clearScreen()
	cyclone := `
                   _                   
  ____ _   _  ____| | ___  ____  _____ 
 / ___) | | |/ ___) |/ _ \|  _ \| ___ |
( (___| |_| ( (___| | |_| | | | | ____|
 \____)\__  |\____)\_)___/|_| |_|_____)
      (____/                           
`
	fmt.Println(cyclone)
	time.Sleep(1 * time.Second)
	clearScreen()
}
