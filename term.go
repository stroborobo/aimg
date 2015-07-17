package aimg

import "fmt"

// cursorUp returns the escape code to set the cursor n lines before the
// current position.
func cursorUp(count int) string {
	return fmt.Sprintf("\033[%dA", count)
}

// newLine returns a string containing an escape code to reset all colors and a
// space. This is needed to prevent artifacts after the end of this line if the
// terminal gets resized.
func newLine() string {
	return "\033[0m \n"
}
