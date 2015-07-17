package main

import "fmt"

func cursorUp(count int) string {
	return fmt.Sprintf("\033[%dA", count)
}

func newLine() string {
	// add a space to prevent artifacts after resizing
	return "\033[0m \n"
}
