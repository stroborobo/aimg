package main

import "fmt"

func cursorUp(count int) {
	fmt.Printf("\033[%dA", count)
}

func reset() {
	// add a space to prevent artifacts after resizing
	fmt.Printf("\033[0m ")
}
