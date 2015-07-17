package main

import "fmt"

func cursorUp(count int) {
	fmt.Printf("\033[%dA", count)
}

func newLine() {
	// add a space to prevent artifacts after resizing
	fmt.Println("\033[0m ")
}
